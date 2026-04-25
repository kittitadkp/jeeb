# Jeeb

Personal management app: workouts, study, sleep, finance, calendar notifications.

## Prerequisites

- Docker and Docker Compose installed
- At least 4GB of available RAM (recommended)

## Quick Start

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd jeeb
   ```

2. Start all services:
   ```bash
   docker-compose up --build
   ```

   Or run in background:
   ```bash
   docker-compose up -d --build
   ```

3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Keycloak (Auth): http://localhost:8081
   - MongoDB: localhost:27017

## Services

- **Frontend** (Port 3000): React application with TypeScript
- **Backend** (Port 8080): Go API server with Clean Architecture
- **MongoDB** (Port 27017): Database for storing application data
- **Keycloak** (Port 8081): Authentication and authorization server

## Environment Variables

The application uses default environment variables for development. For production deployment, create appropriate `.env` files:

- `backend/env/.env.docker` - Backend configuration
- Environment variables in `docker-compose.yml` can be overridden

Default credentials:
- MongoDB: `jeeb/jeeb123`
- Keycloak Admin: `admin/admin123`

## Development

### Running Individual Services

```bash
# Start only database and auth
docker-compose up mongodb keycloak

# Start backend only (requires mongodb and keycloak)
docker-compose up backend

# Start frontend only (requires backend)
docker-compose up frontend
```

### Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
```

### Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: This deletes all data)
docker-compose down -v
```

## API Documentation

See `docs/api/endpoints.md` for API endpoint documentation.

## Architecture

The application follows Clean Architecture principles:
- Domain layer: Business entities and rules
- Use case layer: Application logic
- Adapter layer: External interfaces (HTTP, MongoDB, etc.)

## Troubleshooting

See `docs/troubleshooting/` for common issues and solutions.