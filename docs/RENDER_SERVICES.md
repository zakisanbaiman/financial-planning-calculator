# Render.com Services Documentation

## Current Production Services

The following services are currently deployed and actively serving production traffic:

### Backend API
- **Service Name**: `financial-planning-backend-5n5o`
- **Type**: Web Service
- **Status**: Live ✅
- **URL**: Connected via `NEXT_PUBLIC_API_URL` environment variable
- **Last Deployed**: January 28, 2026

### Frontend Web
- **Service Name**: `financial-planning-frontend-5n5o`  
- **Type**: Web Service
- **Status**: Live ✅
- **URL**: Production URL from Render dashboard
- **Last Deployed**: January 28, 2026

### Database
- **Service Name**: `financial-planning-db`
- **Type**: PostgreSQL Database
- **Plan**: Free tier
- **Region**: Oregon

## Deprecated/Old Services

The following services are deprecated and should be ignored by monitoring:

- `financial-planning-backend` - Replaced by `financial-planning-backend-5n5o`
- `financial-planning-frontend` - Replaced by `financial-planning-frontend-5n5o`
- `financial-planning-frontend-c1zz` - Old test deployment
- `financial-planning-calculator` - Old combined service

## Service Naming Convention

Render.com appends a unique suffix (e.g., `-5n5o`) to service names to ensure uniqueness across the platform. The current production services use the `-5n5o` suffix.

## Monitoring Configuration

The deployment monitoring script (`scripts/check-render-deployments.js`) is configured to:
- Ignore deployments older than 14 days
- Report only on recent deployment failures
- Filter out old/deprecated services automatically

## Deployment Configuration

Services are defined in `render.yaml` at the repository root. This file specifies:
- Service names (without suffix)
- Docker build configuration
- Environment variables
- Health check endpoints
- Build filters for monorepo support

## Troubleshooting

### Old Services Showing as Failed

If monitoring reports show old services as failed:
1. Check the deployment date - if it's more than 14 days old, it will be automatically filtered
2. Verify current production services are healthy (look for `-5n5o` suffix)
3. Old services can be safely ignored if newer services are running

### Actual Production Failures

If a current production service (`-5n5o` suffix) shows as failed:
1. Check the Render dashboard for detailed logs
2. Review recent code changes that triggered the deployment
3. Check GitHub Actions workflow logs for build errors
4. Verify environment variables are correctly configured

## Links

- [Render Dashboard](https://dashboard.render.com)
- [Render.yaml Configuration](../render.yaml)
- [Monitoring Script](../scripts/check-render-deployments.js)
