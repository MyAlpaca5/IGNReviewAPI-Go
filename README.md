# IGNReviewAPI-Go

## Prerequisite
[Go 1.18](https://go.dev/doc/)

[PostgreSQL 14](https://www.postgresql.org/docs/14/index.html)

[golang-migrate v4](https://github.com/golang-migrate/migrate)

[just](https://github.com/casey/just)
 
## API Endpoints
| Method | URL Pattern                  | Action                 
| ------ | ---------------------------- | ------------------------------- 
| GET    | `/healthcheck`               | Show application general application info
| GET    | `/api/admin/metrics`         | Show application application metrics
| GET    | `/api/reviews`               | Show the details of all reviews
| POST   | `/api/reviews`               | Create a new review record
| GET    | `/api/reviews/:id`           | Show the details of a specific review   
| PATCH  | `/api/reviews/:id`           | Update the details of a specific review   
| DELETE | `/api/reviews/:id`           | Delete a specific review
| POST   | `/api/users`                 | Create a new user record
| POST   | `/api/tokens/authentication` | Generate a new access token


## Comments
1. pq: SSL is not enabled on the server -> sslmode=disable, used only for development!!
2. token-based authentication using [PASETO](https://github.com/o1egl/paseto)