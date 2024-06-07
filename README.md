# 🍿 Movie Rating System

The original documentation containing the full challenge description is [here](README_CHALLENGE.md).

## 🚀 System Architecture in a high level

![MRS Architecture](architecture.png)

## 📚 Documentations

Each system component (service or client) is documented in its respective repository:

- ### 👤 User Service [repository](services/user)
- ### ⭐ Rating Service [repository](services/rating)
- ### 🎥 Movie Service [repository](services/movie)
- ### 🔑 Authentication Service [repository](services/authentication)
- ### 💻 Client [repository](client)

## 🏃 How to run

In order to run the system locally, there are two simple alternatives. Note that both options require [Docker](https://www.docker.com/):

**1. From this directory, run:**
```bash
make compose
```
This command will spin up all the containerized components.

**2. See the [documentations](#-documentations) to learn how to start each service or client independently. For the databases, either you run your own instances of [MongoDB](https://www.mongodb.com/) and [Redis](https://redis.io/) or just run:**
```bash
make run-db
```
This command will spin up all the 3 MongoDB and one Redis containers.

### Database migration

To populate the database with 1M+ movies data, run:
```bash
cd services/movie && make migrate
```
The dataset used in this task is provided by [TMDB](https://www.themoviedb.org/) under [ODC Attribution License (ODC-By)](https://opendatacommons.org/licenses/by/1-0/index.html).

## ✅ Requirements

- [x] The backend should expose RESTful endpoints to handle user input and
  return movie ratings.
- [x] The system should store data in a database. You can use any existing
  dataset or API to populate the initial database.
- [x] Implement user endpoints to create and view user information.
- [x] Implement movie endpoints to create and view movie information.
- [x] Implement a rating system to rate the entertainment value of a movie.
- [x] Implement a basic profile where users can view their rated movies.
- [x] Include unit tests to ensure the reliability of your code.
- [x] Ensure proper error handling and validation of user inputs.

## ✨ Bonus Points

- [x] Implement authentication and authorization mechanisms for users.
- [x] Provide documentation for your API endpoints using tools like Swagger.
- [x] Implement logging to record errors and debug information.
- [x] Implement caching mechanisms to improve the rating system's performance.
- [ ] Implement CI/CD quality gates.

## 🌄 Future Improvements

  - Update/delete/list User endpoints
  - Update/delete/list Movie endpoints
  - Monitoring/observability
  - GraphQL or an API Gateway/Routing Service and communication through gRPC
  - Recommendation Service
  - Production-ready WebApp
