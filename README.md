# TamagoAM Backend

## GraphQL endpoints
- POST /graphql
- GET /playground
- GET /health

## Configuration
Copy .env.example to .env and set DB_USER/DB_PASS.

## Migrations
Set MIGRATE_ON_START=true to apply migrations on startup.

## Example query
```
query {
  users {
    id
    name
    userName
    email
    creationDate
  }
}
```
