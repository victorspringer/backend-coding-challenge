# Movie Rating System - 🎥 Movie Service

This Golang service handles Movies creation and visualization for the Movie Rating System backend.

## Requirements

- Go (v1.18 or higher)
- MongoDB

## Setup

1. Install dependencies using `go mod tidy`.
2. Set up your MongoDB instance and update the connection details in the [configs/development.toml](configs/development.toml) file. Or just run the `make run-db` command in the root directory of this monorepository.
3. Run the service using `make run`.
4. Optional: run `make migrate` and populate the database with a large dataset.
5. To run the unit tests, use `make test`.

## Features

- Light speed in-memory cache.
- 1M+ movies dataset.

## Technologies Used

- Golang
- MongoDB

## Endpoints

- You can check, try out and see the models and status codes for every endpoint in the Swagger documentation. With the service running, it is accessible via [http://localhost:8083/docs](http://localhost:8083/docs)
