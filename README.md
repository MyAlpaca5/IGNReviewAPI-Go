# IGN Review API
An application to simulate a simple REST API service using Gin Framework and PostgreSQL.

## Prerequisite
[Go 1.18](https://go.dev/doc/)

[PostgreSQL 14](https://www.postgresql.org/docs/14/index.html)

[golang-migrate v4](https://github.com/golang-migrate/migrate)

[just](https://github.com/casey/just)

## Quick Start
1. Set up PostgreSQL and create an empty database. Then update the environment variables in the `.env` file to reflect your setup
    - info on [Installation and Create Database](https://www.postgresql.org/docs/current/tutorial-start.html) and info on [Connection String](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)
2. Download go dependencies, run `go mod download`
3. Create tables, run `just migration_up`
    - you can remove all tables by running `just migration_down`
    - you can also populate the tables with some fake data by running `just populate_fake_data`, it will create 300+ review records and two users (["username":"admin", "password":"pa55w0rd", "role":"admin"], ["username":"simple", "password":"simplepwd", "role":"simple"])
4. Start the service, run `just run_dev`
    - the application is meant to run in developing environment, for production, need to implement a better way for handling configuration files
 
## API Endpoints
| Method | URL Pattern                  | Required Role     | Description                 
| ------ | ---------------------------- | ----------------- | ------------------------------- 
| GET    | `/healthcheck`               | `None`            | Show application general application info
| POST   | `/api/tokens/authentication` | `None`            | Generate a new access token
| POST   | `/api/users`                 | `None`            | Create a new user record
| GET    | `/api/reviews`               | at least `Simple` | Show the details of all reviews
| POST   | `/api/reviews`               | at least `Simple` | Create a new review record
| GET    | `/api/reviews/:id`           | at least `Simple` | Show the details of a specific review   
| PATCH  | `/api/reviews/:id`           | at least `Simple` | Update the details of a specific review   
| DELETE | `/api/reviews/:id`           | at least `Simple` | Delete a specific review
| GET    | `/api/admin/metrics`         | `Admin`           | Show application application metrics

**Accepted Query Paramaters**
| Method | URL Pattern                  | Required                             | Optional                 
| ------ | ---------------------------- | ------------------------------------ | ------------------------------- 
| GET    | `/healthcheck`               | `None`                               | `None`
| POST   | `/api/tokens/authentication` | `username`, `password`               | `None`
| GET    | `/api/reviews`               | `None`                               | `name`, `score_min`, `order`, `page`, `page_size`, `genres`
| POST   | `/api/reviews`               | `name`, `review_url`, `review_score` | `description`, `media_type`, `genre_list`, `creator_list`
| GET    | `/api/reviews/:id`           | `None`                               | `None`
| PATCH  | `/api/reviews/:id`           | `None`                               | `name`, `review_url`, `review_score`, `description`, `media_type`, `genre_list`, `creator_list`
| DELETE | `/api/reviews/:id`           | `None`                               | `None`
| POST   | `/api/users`                 | `username`, `password`               | `email`
| GET    | `/api/admin/metrics`         | `None`                               | `None`

## Comments
1. solution for error `pq: SSL is not enabled on the server` -> set `sslmode=disable`, should only be used only for development!!
2. token-based authentication using [PASETO](https://github.com/o1egl/paseto)
