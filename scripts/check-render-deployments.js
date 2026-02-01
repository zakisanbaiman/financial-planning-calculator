#!/usr/bin/env node

/**
 * Check Render.com Deployments
 * 
 * This script checks the status of recent deployments on Render.com
 * and reports any errors or failures.
 * 
 * Usage:
 *   RENDER_API_KEY=xxx node scripts/check-render-deployments.js
 */

const https = require('https');

const RENDER_API_KEY = process.env.RENDER_API_KEY;
const RENDER_OWNER_ID = process.env.RENDER_OWNER_ID;

if (!RENDER_API_KEY) {
  console.error('❌ Error: RENDER_API_KEY environment variable is required');
  process.exit(1);
}

/**
 * Make HTTPS request to Render API
 */
function makeRenderAPIRequest(path, method = 'GET') {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'api.render.com',
      port: 443,
      path: path,
      method: method,
      headers: {
        'Authorization': `Bearer ${RENDER_API_KEY}`,
        'Accept': 'application/json',
      },
    };

    const req = https.request(options, (res) => {
      let data = '';

      res.on('data', (chunk) => {
        data += chunk;
      });

      res.on('end', () => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          try {
            resolve(JSON.parse(data));
          } catch (e) {
            resolve(data);
          }
        } else {
          reject(new Error(`API request failed with status ${res.statusCode}: ${data}`));
        }
      });
    });

    req.on('error', (error) => {
      reject(error);
    });

    req.end();
  });
}

/**
 * List all services
 */
async function listServices() {
  const params = RENDER_OWNER_ID ? `?ownerId=${RENDER_OWNER_ID}` : '';
  const response = await makeRenderAPIRequest(`/v1/services${params}`);
  return response;
}

/**
 * List deployments for a service
 */
async function listDeployments(serviceId, limit = 5) {
  const response = await makeRenderAPIRequest(`/v1/services/${serviceId}/deploys?limit=${limit}`);
  return response;
}

/**
 * Get deployment logs
 */
async function getDeploymentLogs(serviceId, deployId) {
  const logs = await makeRenderAPIRequest(`/v1/services/${serviceId}/deploys/${deployId}/logs`);
  return logs;
}

/**
 * Detect errors in logs
 */
function detectErrors(logs) {
  const errors = [];
  const errorPatterns = [
    {
      pattern: /error|failed|failure/i,
      type: 'general_error',
      severity: 'high',
    },
    {
      pattern: /npm ERR!/i,
      type: 'npm_error',
      severity: 'high',
    },
    {
      pattern: /fatal/i,
      type: 'fatal_error',
      severity: 'critical',
    },
    {
      pattern: /cannot find module|module not found/i,
      type: 'missing_dependency',
      severity: 'high',
    },
    {
      pattern: /ECONNREFUSED|ETIMEDOUT/i,
      type: 'connection_error',
      severity: 'high',
    },
    {
      pattern: /syntax error/i,
      type: 'syntax_error',
      severity: 'high',
    },
    {
      pattern: /out of memory|OOM/i,
      type: 'memory_error',
      severity: 'critical',
    },
  ];

  const lines = logs.split('\n');
  lines.forEach((line, index) => {
    errorPatterns.forEach((pattern) => {
      if (pattern.pattern.test(line)) {
        errors.push({
          line: index + 1,
          content: line.trim(),
          type: pattern.type,
          severity: pattern.severity,
        });
      }
    });
  });

  return errors;
}

/**
 * Main function
 */
async function main() {
  try {
    console.log('🔍 Checking Render.com deployments...\n');

    // List all services
    const services = await listServices();
    console.log(`📦 Found ${services.length} services\n`);

    let hasErrors = false;
    const report = [];
    
    // Filter threshold: Only check deployments from last 14 days
    const FOURTEEN_DAYS_AGO = Date.now() - (14 * 24 * 60 * 60 * 1000);

    // Check each service
    for (const serviceData of services) {
      const service = serviceData.service;
      console.log(`\n📋 Service: ${service.name}`);
      console.log(`   Type: ${service.type}`);
      console.log(`   Status: ${service.serviceDetails?.state || 'unknown'}`);

      // Get recent deployments
      try {
        const deploys = await listDeployments(service.id, 3);
        
        if (deploys.length === 0) {
          console.log('   ℹ️  No deployments found');
          continue;
        }

        console.log(`   📊 Recent deployments: ${deploys.length}`);

        // Check latest deployment
        const latestDeploy = deploys[0].deploy;
        const deployDate = new Date(latestDeploy.createdAt).getTime();
        console.log(`\n   Latest Deploy:`);
        console.log(`   - ID: ${latestDeploy.id}`);
        console.log(`   - Status: ${latestDeploy.status}`);
        console.log(`   - Created: ${new Date(latestDeploy.createdAt).toLocaleString()}`);

        // Skip old deployments - they're likely deprecated services
        // Only check deployments from the last 14 days
        if (deployDate < FOURTEEN_DAYS_AGO) {
          console.log(`   ℹ️  Deployment is older than 14 days - skipping check`);
          continue;
        }

        if (latestDeploy.status === 'build_failed' || latestDeploy.status === 'deactivated') {
          hasErrors = true;
          console.log(`   ⚠️  DEPLOYMENT FAILED!`);

          // Get logs for failed deployment
          try {
            const logs = await getDeploymentLogs(service.id, latestDeploy.id);
            const errors = detectErrors(logs);

            if (errors.length > 0) {
              console.log(`\n   🔴 Detected ${errors.length} errors:`);
              
              // Show first 5 errors
              errors.slice(0, 5).forEach((error) => {
                console.log(`      [${error.severity}] ${error.type}`);
                console.log(`      Line ${error.line}: ${error.content.substring(0, 100)}...`);
              });

              if (errors.length > 5) {
                console.log(`      ... and ${errors.length - 5} more errors`);
              }

              report.push({
                service: service.name,
                deployId: latestDeploy.id,
                status: latestDeploy.status,
                errors: errors,
              });
            }
          } catch (logError) {
            console.log(`   ⚠️  Could not fetch logs: ${logError.message}`);
          }
        } else if (latestDeploy.status === 'live') {
          console.log(`   ✅ Deployment is live and healthy`);
        } else {
          console.log(`   ⏳ Deployment in progress...`);
        }
      } catch (deployError) {
        console.log(`   ⚠️  Could not fetch deployments: ${deployError.message}`);
      }
    }

    // Print summary
    console.log('\n' + '='.repeat(60));
    console.log('📊 SUMMARY');
    console.log('='.repeat(60));

    if (hasErrors) {
      console.log('❌ Some deployments have errors or failures');
      console.log('\n🔧 Recommended actions:');
      console.log('1. Check the error logs above');
      console.log('2. Fix the identified issues in your code');
      console.log('3. Push a new commit to trigger redeployment');
      console.log('4. Use MCP tools to get detailed error analysis');
      
      if (report.length > 0) {
        console.log('\n📝 Error Report:');
        console.log(JSON.stringify(report, null, 2));
      }

      process.exit(1);
    } else {
      console.log('✅ All deployments are healthy');
      process.exit(0);
    }
  } catch (error) {
    console.error('❌ Fatal error:', error.message);
    process.exit(1);
  }
}

main().catch((error) => {
  console.error('❌ Unhandled error:', error);
  process.exit(1);
});
