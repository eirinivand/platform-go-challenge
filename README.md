# Asset API

This is just another API, or is it?

## Introduction

This API is designed to allow users to interact with Asset data (stored in a MongoDB).

### Roles

The users in this concept have roles assigned to them.
The roles determined which endpoints they can use from the API.

There are two roles currently existing; `admin` and `user`.

### Asset Types

| Role  | Endpoints                    | Action               |
|:------|:-----------------------------|:---------------------|
| N/A   | api/v1/users/signup          | SignUp               | 
| Admin | api/v1/admin/add-charts-bulk | Add a bulk of Charts |
|       |                              |                      |

## Usage

### Prerequisites

There are two ways to run this application;

1. Docker Compose (includes a MongoDB)
2. Go Run (along with a MongoDB).

Specifically:

1. For both approaches we need an `.env` file in the root directory of the project
   (same level with this `README.md` and `main.go` file). Please don't commit this file.
    ```ini
    MONGODB=mongodb://localhost:27017
    DB_NAME=favourites
    JWT_SECRET_KEY=SOOOOOO_SECRET_KEY
   ```
   where:
    * `MONGODB` is the URI of the MongoDB
    * `DB_NAME` is the name of the database to be used by the application
    * `JWT_SECRET_KEY` is the Secret Key used to sign the Json Web Token
      used during authentication and authorization
2. A Postman installation (see [official instructions here](https://www.postman.com/downloads/))
    1. Import the [GoFavourites.postman_environment.json](GoFavourites.postman_environment.json) and then
    2. import [GoFavourites.postman_collection.json](GoFavourites.postman_collection.json)
3. Then we can either:
    1. (Recommended) Run with Docker Compose
        1. Have Docker and Docker Compose installed
           (see [official instructions here](https://docs.docker.com/compose/install/)).
    2. Run using Go Run for which we need to
        1. have set `GOROOT`, `GOPATH`
        2. Have a working MongoDB running (and the URI in `.env` file)
        3. Then open a terminal, navigate in this directory
        4. Run `go run main.go`

## Documentation

### Postman Collection Usage

First we need a user, thus we should create one.

## Other Info

Based on [Instructions of the challenge](Instructions.md)
