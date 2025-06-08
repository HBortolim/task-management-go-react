# Task Management API

A RESTful API for managing tasks and goals built with Go, Gin, and MongoDB.

## Requirements

- Go 1.x
- MongoDB
- Git

## Installation and Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd task-management-go-react/backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
- Copy `.env.example` to `.env` (or use the existing .env file)
- Update the values in `.env` as needed

4. Start MongoDB:
- Make sure MongoDB is running on your system
- Default connection URI: `mongodb://localhost:27017`

5. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on http://localhost:8080

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and get JWT token

### Goals

All goal endpoints require authentication (JWT token in Authorization header)

- `GET /api/goals` - Get all goals for the logged-in user
- `GET /api/goals/:id` - Get a specific goal
- `POST /api/goals` - Create a new goal
- `PUT /api/goals/:id` - Update a goal
- `DELETE /api/goals/:id` - Delete a goal

### Health Check

- `GET /health` - Check if the API is running

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Main application entry point
├── configs/
│   └── config.go            # Configuration handling
├── internal/
│   ├── db/
│   │   └── mongodb.go       # MongoDB connection
│   ├── handlers/
│   │   ├── auth.go          # Authentication handlers
│   │   ├── goal.go          # Goal CRUD handlers
│   │   └── routes.go        # Route setup
│   ├── middleware/
│   │   └── auth.go          # JWT authentication middleware
│   └── models/
│       ├── user.go          # User model
│       └── goal.go          # Goal and SubTask models
├── .env                     # Environment variables
└── README.md                # This file
```

## Docker (Future Implementation)

Docker and Kubernetes support will be added in future versions.

## License

MIT 