# Chirpy - A Twitter-Like API Service

## Overview
Chirpy is a backend API service built with Go that provides core functionality similar to Twitter/X. Users can create accounts, post short messages (chirps), and interact with other users' content. The service includes JWT-based authentication, database persistence, and various API endpoints to manage users and chirps.

## Features

### Authentication System
- User registration and login
- JWT-based access tokens
- Refresh token system for prolonged sessions
- Token revocation

### User Management
- Create user accounts
- Update user profiles
- Premium (Chirpy Red) subscription support

### Chirp Functionality
- Create chirps (140 character limit)
- List chirps with sorting and filtering
- Delete chirps (author-only)
- Profanity filtering

### Admin Features
- Usage metrics
- Database reset (development mode only)

### Webhook Integration
- Support for Polka payment service integration
- User subscription status updates

## Technical Stack
- **Backend:** Go
- **Database:** PostgreSQL
- **ORM:** SQLC for type-safe SQL
- **Migration:** Goose
- **Authentication:** JWT 
- **Configuration:** Environment variables via `godotenv`



### More About Development Tools

#### SQLC Commands

Installing sqlc CLI tools with `go install`
> go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

SQLC generates type-safe Go code from SQL queries.
Before using SQLC, you need to create a configuration file named `sqlc.yaml` in your project root:

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

```bash
# Generate Go code from SQL queries
sqlc generate

# Verify SQL queries without generating code
sqlc vet
```

Example SQLC query definition:
```sql
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
```
Tutorial : 
https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

#### Gooose Migrations


Installing goose CLI tools with `go install`
> go install github.com/pressly/goose/v3/cmd/goose@latest

Goose handles database schema migrations.

```bash
# Create a new migration
goose -dir sql/schema create add_users_table sql

# Apply migrations
goose -dir sql/schema postgres "postgresql://user:password@localhost:5432/chirpy?sslmode=disable" up

# Roll back the most recent migration
goose -dir sql/schema postgres "postgresql://user:password@localhost:5432/chirpy?sslmode=disable" down

# Check migration status
goose -dir sql/schema postgres "postgresql://user:password@localhost:5432/chirpy?sslmode=disable" status
```

Example migration file:
```sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
```


## API Endpoints

### Authentication
| Method | Endpoint       | Description          |
| ------ | -------------- | -------------------- |
| POST   | `/api/users`   | Create a new user    |
| POST   | `/api/login`   | Login and get tokens |
| POST   | `/api/refresh` | Refresh access token |
| POST   | `/api/revoke`  | Revoke refresh token |

### Users
| Method | Endpoint     | Description         |
| ------ | ------------ | ------------------- |
| PUT    | `/api/users` | Update user profile |

### Chirps
| Method | Endpoint                | Description                              |
| ------ | ----------------------- | ---------------------------------------- |
| POST   | `/api/chirps`           | Create a new chirp                       |
| GET    | `/api/chirps`           | Get all chirps (with optional filtering) |
| GET    | `/api/chirps/{chirpID}` | Get a specific chirp                     |
| DELETE | `/api/chirps/{chirpID}` | Delete a chirp                           |

### Admin
| Method | Endpoint         | Description                        |
| ------ | ---------------- | ---------------------------------- |
| GET    | `/admin/metrics` | View API usage metrics             |
| POST   | `/admin/reset`   | Reset the database (dev mode only) |

### Webhooks
| Method | Endpoint              | Description                |
| ------ | --------------------- | -------------------------- |
| POST   | `/api/polka/webhooks` | Handle subscription events |

## Installation and Setup

### Prerequisites
- Go 1.24 or later
- PostgreSQL
- Environment variables (see below)

### Environment Variables
Create a `.env` file in the project root with the following:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=chirpy
JWT_SECRET=your_jwt_secret
```


## Note
This project was built as part of the Boot.dev backend programming curriculum, designed to provide hands-on experience with building a RESTful API service in Go.