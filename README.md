# IGNReviewAPI-Go

## Prerequisite
Go 1.18
PostgreSQL
 
## API Endpoints
| Method | URL Pattern                  | Action                 
| ------ | ---------------------------- | ------------------------------- 
| GET    | `/api/healthcheck`           | Show application health and version information
| GET    | `/debug/vars`                | Show application metrics
| GET    | `/api/reviews`               | Show the summary of all reviews
| POST   | `/api/reviews`               | Create a new review entry
| GET    | `/api/reviews/:id`           | Show the details of a specific review   
| PATCH  | `/api/reviews/:id`           | Update the details of a specific review   
| DELETE | `/api/reviews/:id`           | Delete a specific review
| POST   | `/api/users`                 | Create a new user entry
| PUT    | `/api/users/activate`        | Activate a specific user
| PUT    | `/api/users/password`        | Update the password of a specific user
| POST   | `/api/tokens/authentication` | Generate a new authentication token
| POST   | `/api/tokens/password-reset` | Generate a new password-reset token
