# Load all variables from .env
ifneq (,$(wildcard .env))
    include .env
    export
endif

MIGRATIONS_PATH= ./migrations

.PHONY : migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir ${MIGRATIONS_PATH} $(NAME)

.PHONY : migrate-up
migrate-up:
	@migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" up

.PHONY : migrate-goto
migrate-goto:
	@migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" goto $(VERSION)

.PHONY : migrate-down
migrate-down:
	@migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" down

.PHONY : migrate-roll
migrate-roll:
	@migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" down $(VERSION)

.PHONY : migrate-force
migrate-force:
	@migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" force $(VERSION)

.PHONY : gen-docs
gen-docs:
	@echo "Generating API documentation..."
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
	@echo "API documentation generated successfully."

.PHONY :  test
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests completed."