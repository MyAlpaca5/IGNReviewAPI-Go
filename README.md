# IGNReviewAPI-Go

## Prerequisite
[Go 1.18](https://go.dev/doc/)

[PostgreSQL 14](https://www.postgresql.org/docs/14/index.html)

[golang-migrate v4](https://github.com/golang-migrate/migrate)

[make](https://www.gnu.org/software/make/) or [just](https://github.com/casey/just)
 
## API Endpoints
| Method | URL Pattern                  | Action                 
| ------ | ---------------------------- | ------------------------------- 
| GET    | `/healthcheck`               | Show application health and version information
| GET    | `/api/admin/metrics`         | Show application metrics
| GET    | `/api/reviews`               | Show the summary of all reviews
| POST   | `/api/reviews`               | Create a new review entry
| GET    | `/api/reviews/:id`           | Show the details of a specific review   
| PATCH  | `/api/reviews/:id`           | Update the details of a specific review   
| DELETE | `/api/reviews/:id`           | Delete a specific review
| POST   | `/api/users`                 | Create a new user entry
| PUT    | `/api/users/password`        | Update the password of a specific user
| POST   | `/api/tokens/authentication` | Generate a new authentication token

