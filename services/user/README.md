# Movie Rating System - 👤 User Service

This Golang service handles Users creation and visualization for the Movie Rating System backend.

## Requirements

- Go (v1.18 or higher)
- MongoDB

## Setup

1. Install dependencies using `go mod tidy`.
2. Set up your MongoDB instance and update the connection details in the [configs/development.toml](configs/development.toml) file. Or just run the `make run-db` command in the root directory of this monorepository.
3. Run the service using `make run`.
4. To run the unit tests, use `make test`.

## Technologies Used

- Golang
- MongoDB

## Endpoints

- You can check, try out and see the models and status codes for every endpoint in the Swagger documentation. With the service running, it is accessible via [http://localhost:8081/docs](http://localhost:8081/docs)
