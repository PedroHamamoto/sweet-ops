api-start:
	go run cmd/main.go
vet:
	go vet ./...

MIGRATE_BIN ?= migrate
MIGRATIONS_DIR ?= ./migrations
DB_URL ?= postgres://sweetops:sweetops_dev_password@localhost:7777/sweetops_db?sslmode=disable

migrate-up:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1

migrate-version:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) version

migrate-force:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database $(DB_URL) force $(VERSION)

migrate-create:
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) $(NAME)

DB_CONTAINER ?= sweetops-postgres
DB_USER ?= sweetops
DB_NAME ?= sweetops_db
DB_QUERY ?= SELECT 1;

db-query:
	docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "$(DB_QUERY)"