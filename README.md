# Sweet Ops

A bakery management platform 

---

## Requirements

- **[Go](https://go.dev/)** (version 1.25+)
- **[Docker & Docker Compose](https://www.docker.com/)**
- **[golang-migrate](https://github.com/golang-migrate/migrate)**

---

## Getting Started

### 1. Set up environment variables

Copy the example file and adjust values as needed:

```bash
cp .env.example .env
```

### 2. Spin up the database

```bash
docker compose up -d
```

### 3. Run the migrations

```bash
make migrate-up
```

### 4. Start the API

```bash
make api-start
```

The API will be running at http://localhost:8080

---

## API

| Method | Endpoint     | Description         |
|--------|--------------|---------------------|
| POST   | /api/users   | Create a new user   |
| POST   | /api/login   | Authenticate a user |

---

## Development

### Creating migrations

To create a new set of migration files (up and down) with the current timestamp:

```bash
make migrate-create NAME=create_example_table
```

### Rolling back a migration

```bash
make migrate-down
```

### Vet

```bash
make vet
```
