# Golang Social Project

This project is a social application built with Go and SQL.

## Prerequisites

- Go version 1.16 or later
- PostgreSQL

## Setup

1. Clone the repository:

```bash
git clone https://github.com/minggitprakasa/golang_social.git
cd golang_social
```

2. Copy the .env.example file and create a new .env file:

```bash
cp .env.example .env
```

3. Install the dependencies:

```bash
go mod download
```

4. Migration database:

```bash
migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" -path db/migrations up
```

5. Run the application:

```bash
go run ./cmd/main.go
```

