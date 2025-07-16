# Tải các biến từ file .env và export chúng
include .env
export

# Định nghĩa các biến để tái sử dụng
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
POSTGRES_CONTAINER_NAME := ecommerce_db 

# .PHONY để khai báo các target không phải là file
.PHONY: postgres createdb dropdb migrateup migratedown sqlc server compose-up compose-down

# ==================== Docker Commands ====================
compose-up:
	docker compose up -d

compose-down:
	docker compose down

postgres:
	docker compose up -d postgres 
# ==================== Database Commands ====================
createdb:
	docker exec -it $(POSTGRES_CONTAINER_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it $(POSTGRES_CONTAINER_NAME) dropdb $(DB_NAME)

migrateup:
	goose -dir "db/migration" postgres "$(DB_URL)" up

migratedown:
	goose -dir "db/migration" postgres "$(DB_URL)" down -v

# ==================== App Commands ====================
sqlc:
	sqlc generate

server:
	go run cmd/app/main.go

# Lệnh mặc định khi gõ `make`
default: server