# GEMINI.md

This file provides a comprehensive overview of the financial-planning-calculator project, designed to be used as a context for AI-powered development assistants like Gemini.

## Project Overview

The financial-planning-calculator is a web application designed to help users visualize their future asset growth and retirement financial planning. Users can input their current income, expenses, and savings to calculate future asset transitions, retirement funds, and emergency preparedness, enabling them to create a secure financial plan.

**Technologies:**

*   **Frontend:** Next.js 14, TypeScript, Tailwind CSS, Chart.js, React Hook Form, Zod
*   **Backend:** Go, Echo Framework, PostgreSQL
*   **Containerization:** Docker

**Architecture:**

The project follows a standard monorepo structure with two main components:

*   `frontend/`: A Next.js application for the user interface.
*   `backend/`: A Go application providing a RESTful API.

The frontend and backend are designed to be developed and run independently, but are orchestrated using Docker for a unified development environment.

## Building and Running

This project supports both Docker and local development environments. Docker is the recommended approach.

### Docker Development Environment (Recommended)

**Prerequisites:**

*   Docker
*   Docker Compose
*   `make`

**Commands:**

*   `make dev-setup`: First-time setup (builds containers, starts the database, runs migrations, and seeds data).
*   `make up`: Start the application.
*   `make down`: Stop the application.
*   `make test`: Run backend tests.
*   `make lint`: Run backend linter.
*   `make shell-api`: Access the backend container shell.
*   `make shell-db`: Access the database container shell.
*   `make logs`: View application logs.

### Local Development Environment

**Prerequisites:**

*   Node.js 18.20.0+
*   Go 1.24.0+
*   PostgreSQL 13+

**Frontend:**

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

**Backend:**

```bash
cd backend
go mod tidy
cp .env.example .env
go run main.go
```

## Development Conventions

### Frontend

*   **Linting:** `npm run lint` (ESLint)
*   **Type Checking:** `npm run type-check` (TypeScript)
*   **Package Management:** `npm`

### Backend

*   **Testing:** `go test ./...`
*   **Dependency Management:** Go Modules (`go mod tidy`)
*   **API Documentation:** Swagger (available at `http://localhost:8080/swagger/index.html`)
*   **Performance Profiling:** `pprof` is enabled in the development environment.

### Version Control

This project supports `direnv`, `goenv`, and `asdf` for managing Go and Node.js versions. `direnv` is the recommended tool for automatically using the Go version specified in the Docker container.
