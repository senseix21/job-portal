
# Job Portal Backend

A backend application for managing job portal functionalities built using Go and the Echo framework. This backend provides APIs for managing users, jobs, and other job portal-related features.
live url: https://job-portal-frontend-pink.vercel.app/

## Table of Contents

- [Installation](#installation)
- [Project Structure](#project-structure)
- [Setup](#setup)
- [API Documentation](#api-documentation)
- [Middleware & Security](#middleware--security)
- [License](#license)

## Installation

To get started with this project, clone the repository and follow the steps below:

### Prerequisites

- [Go](https://golang.org/dl/) version 1.18 or higher
- MongoDB instance (can be local or remote)
- Echo framework (via Go Modules)

### Steps

1. Clone this repository:

   ```bash
   git clone https://github.com/your-username/job-portal-backend.git
   ```

2. Change to the project directory:

   ```bash
   cd job-portal-backend
   ```

3. Install the required dependencies using Go Modules:

   ```bash
   go mod tidy
   ```

4. Set up environment variables for database connection and other services. You can create a `.env` file in the root directory of the project and add necessary configurations (such as MongoDB URI).

5. Run the application:

   ```bash
   go run main.go
   ```

6. The server will start on `http://localhost:8080`.

## Project Structure

```
job-portal-backend/
├── config/                  # Database configuration and initialization
├── controllers/             # Controllers for handling user and job logic
├── middlewares/             # Custom middlewares (e.g., error handler, validation)
├── routers/                 # Route definitions
├── services/                # Services for business logic
├── main.go                  # Main application file
├── go.mod                   # Go module file
├── go.sum                   # Go module checksum file
└── README.md                # This file
```

### Key Directories

- **config**: Handles database connection and configuration.
- **controllers**: Contains logic for handling requests and interacting with services.
- **middlewares**: Custom middlewares, such as custom error handling, validation, and rate limiting.
- **routers**: Defines API routes for users and jobs.
- **services**: Contains business logic related to jobs, users, and other operations.

## Setup

### Environment Variables

Create a `.env` file in the root directory and add the following variables:

```ini
MONGO_URI=mongodb://localhost:27017
```

Replace `localhost` with your MongoDB server URL if you're using a remote instance.

## API Documentation

### User Routes

- **POST `/users/register`** - Register a new user.
- **POST `/users/login`** - Login a user.
- **GET `/users/:id`** - Fetch user details by ID.
  
### Job Routes

- **GET `/jobs`** - List all jobs.
- **POST `/jobs`** - Create a new job posting.
- **GET `/jobs/:id`** - Get details of a specific job by its ID.
- **PUT `/jobs/:id`** - Update a job posting by its ID.
- **DELETE `/jobs/:id`** - Delete a job posting by its ID.

## Middleware & Security

This project utilizes several middleware features to ensure the security, efficiency, and functionality of the API.

### CORS Configuration

Cross-Origin Resource Sharing (CORS) is enabled to allow the frontend to communicate with the backend from specified domains (`https://job-portal-frontend-pink.vercel.app` and `http://localhost:3000`).

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://job-portal-frontend-pink.vercel.app", "http://localhost:3000"},
    AllowMethods: []string{
        http.MethodGet,
        http.MethodPost,
        http.MethodPut,
        http.MethodPatch,
        http.MethodDelete,
        http.MethodOptions,
    },
    AllowHeaders: []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

### Rate Limiting

Rate limiting is enabled to limit requests to 20 requests per second.

```go
e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
```

### Body Limit

Requests are limited to a maximum of 2MB to prevent large payloads from overloading the server.

```go
e.Use(middleware.BodyLimit("2M"))
```

### Custom Error Handler

A custom error handler is used to format error responses consistently:

```go
e.HTTPErrorHandler = middlewares.CustomHTTPErrorHandler
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

This README provides an overview of your backend project, how to set it up, and how to interact with it. If there are specific API routes or features not covered here, you can expand on the documentation as needed. Let me know if you need further adjustments!
