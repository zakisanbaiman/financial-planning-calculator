#!/usr/bin/env node

/**
 * Render.com MCP Server
 * 
 * This MCP server provides tools for monitoring Render.com deployments,
 * detecting errors, and enabling AI assistants to automatically fix issues.
 * 
 * Available tools:
 * - get_deployment_logs: Fetch recent deployment logs
 * - get_service_status: Check service health status
 * - list_recent_deploys: List recent deployments
 * - detect_errors: Analyze logs for common error patterns
 */

const https = require('https');
const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} = require('@modelcontextprotocol/sdk/types.js');

const RENDER_API_KEY = process.env.RENDER_API_KEY;
const RENDER_OWNER_ID = process.env.RENDER_OWNER_ID;

if (!RENDER_API_KEY) {
  console.error('Error: RENDER_API_KEY environment variable is required');
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
  return await makeRenderAPIRequest(`/v1/services${params}`);
}

/**
 * Get service details
 */
async function getService(serviceId) {
  return await makeRenderAPIRequest(`/v1/services/${serviceId}`);
}

/**
 * List deployments for a service
 */
async function listDeployments(serviceId, limit = 10) {
  return await makeRenderAPIRequest(`/v1/services/${serviceId}/deploys?limit=${limit}`);
}

/**
 * Get deployment logs
 */
async function getDeploymentLogs(serviceId, deployId) {
  return await makeRenderAPIRequest(`/v1/services/${serviceId}/deploys/${deployId}/logs`);
}

/**
 * Detect common error patterns in logs
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
      pattern: /port.*already in use/i,
      type: 'port_conflict',
      severity: 'medium',
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
 * Create MCP server
 */
const server = new Server(
  {
    name: 'render-mcp-server',
    version: '1.0.0',
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

/**
 * List available tools
 */
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: 'list_services',
        description: 'List all Render.com services in the account',
        inputSchema: {
          type: 'object',
          properties: {},
        },
      },
      {
        name: 'get_service_status',
        description: 'Get the current status of a specific Render.com service',
        inputSchema: {
          type: 'object',
          properties: {
            serviceName: {
              type: 'string',
              description: 'Name of the service (e.g., financial-planning-backend)',
            },
          },
          required: ['serviceName'],
        },
      },
      {
        name: 'list_recent_deploys',
        description: 'List recent deployments for a service',
        inputSchema: {
          type: 'object',
          properties: {
            serviceName: {
              type: 'string',
              description: 'Name of the service',
            },
            limit: {
              type: 'number',
              description: 'Number of deployments to retrieve (default: 10)',
              default: 10,
            },
          },
          required: ['serviceName'],
        },
      },
      {
        name: 'get_deployment_logs',
        description: 'Fetch deployment logs for a specific deployment',
        inputSchema: {
          type: 'object',
          properties: {
            serviceName: {
              type: 'string',
              description: 'Name of the service',
            },
            deployId: {
              type: 'string',
              description: 'Deployment ID (optional, uses latest if not provided)',
            },
          },
          required: ['serviceName'],
        },
      },
      {
        name: 'detect_errors',
        description: 'Analyze deployment logs to detect common error patterns',
        inputSchema: {
          type: 'object',
          properties: {
            serviceName: {
              type: 'string',
              description: 'Name of the service',
            },
            deployId: {
              type: 'string',
              description: 'Deployment ID (optional, uses latest if not provided)',
            },
          },
          required: ['serviceName'],
        },
      },
    ],
  };
});

/**
 * Handle tool calls
 */
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  try {
    switch (name) {
      case 'list_services': {
        const services = await listServices();
        return {
          content: [
            {
              type: 'text',
              text: JSON.stringify(services, null, 2),
            },
          ],
        };
      }

      case 'get_service_status': {
        const services = await listServices();
        const service = services.find(
          (s) => s.service.name === args.serviceName
        );

        if (!service) {
          return {
            content: [
              {
                type: 'text',
                text: `Service "${args.serviceName}" not found`,
              },
            ],
            isError: true,
          };
        }

        const serviceDetails = await getService(service.service.id);
        return {
          content: [
            {
              type: 'text',
              text: JSON.stringify(serviceDetails, null, 2),
            },
          ],
        };
      }

      case 'list_recent_deploys': {
        const services = await listServices();
        const service = services.find(
          (s) => s.service.name === args.serviceName
        );

        if (!service) {
          return {
            content: [
              {
                type: 'text',
                text: `Service "${args.serviceName}" not found`,
              },
            ],
            isError: true,
          };
        }

        const deploys = await listDeployments(
          service.service.id,
          args.limit || 10
        );
        return {
          content: [
            {
              type: 'text',
              text: JSON.stringify(deploys, null, 2),
            },
          ],
        };
      }

      case 'get_deployment_logs': {
        const services = await listServices();
        const service = services.find(
          (s) => s.service.name === args.serviceName
        );

        if (!service) {
          return {
            content: [
              {
                type: 'text',
                text: `Service "${args.serviceName}" not found`,
              },
            ],
            isError: true,
          };
        }

        let deployId = args.deployId;
        if (!deployId) {
          const deploys = await listDeployments(service.service.id, 1);
          if (deploys.length === 0) {
            return {
              content: [
                {
                  type: 'text',
                  text: 'No deployments found',
                },
              ],
              isError: true,
            };
          }
          deployId = deploys[0].deploy.id;
        }

        const logs = await getDeploymentLogs(service.service.id, deployId);
        return {
          content: [
            {
              type: 'text',
              text: logs,
            },
          ],
        };
      }

      case 'detect_errors': {
        const services = await listServices();
        const service = services.find(
          (s) => s.service.name === args.serviceName
        );

        if (!service) {
          return {
            content: [
              {
                type: 'text',
                text: `Service "${args.serviceName}" not found`,
              },
            ],
            isError: true,
          };
        }

        let deployId = args.deployId;
        if (!deployId) {
          const deploys = await listDeployments(service.service.id, 1);
          if (deploys.length === 0) {
            return {
              content: [
                {
                  type: 'text',
                  text: 'No deployments found',
                },
              ],
              isError: true,
            };
          }
          deployId = deploys[0].deploy.id;
        }

        const logs = await getDeploymentLogs(service.service.id, deployId);
        const errors = detectErrors(logs);

        return {
          content: [
            {
              type: 'text',
              text: JSON.stringify(
                {
                  totalErrors: errors.length,
                  errors: errors,
                  summary: {
                    critical: errors.filter((e) => e.severity === 'critical')
                      .length,
                    high: errors.filter((e) => e.severity === 'high').length,
                    medium: errors.filter((e) => e.severity === 'medium')
                      .length,
                  },
                },
                null,
                2
              ),
            },
          ],
        };
      }

      default:
        return {
          content: [
            {
              type: 'text',
              text: `Unknown tool: ${name}`,
            },
          ],
          isError: true,
        };
    }
  } catch (error) {
    return {
      content: [
        {
          type: 'text',
          text: `Error: ${error.message}`,
        },
      ],
      isError: true,
    };
  }
});

/**
 * Start the server
 */
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error('Render.com MCP server running on stdio');
}

main().catch((error) => {
  console.error('Fatal error in main():', error);
  process.exit(1);
});
