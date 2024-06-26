# Movie Rating System - 🔑 Authentication Service

This Golang service handles user authentication and authorization for the Movie Rating System backend.

## Requirements

- Go (v1.18 or higher)
- Redis

## Setup

1. Install dependencies using `go mod tidy`.
2. Set up your Redis instance and update the connection details in the [configs/development.toml](configs/development.toml) file. Or just run the `make run-db` command in the root directory of this monorepository.
3. Run the service using `make run`.

## Features

- User authentication with JWT tokens.
- Authorization mechanisms for accessing protected endpoints.
- The HTTP middleware that handles authorization accepts both Authorization Bearer header and Cookies (`MRSAccessToken` and `MRSRefreshToken`). Therefore, don't be surprised if your requests to other services are accepted if you provide an invalid Authorization header, for example (as long as you have valid cookies and vice versa).

## Technologies Used

- Golang
- Redis
- JWT
- RSA cryptography

## Endpoints

- You can check, try out and see the models and status codes for every endpoint in the Swagger documentation. With the service running, it is accessible via [http://localhost:8084/docs](http://localhost:8084/docs)


## TO DO

- Add tests
