# Load all variables from .env
ifneq (,$(wildcard .env))
	include .env
	export
endif

MIGRATIONS_PATH=./migrations

# ========================================================================================== #
# HELPERS
# ========================================================================================== #

.PHONY: help confirm

## help: Show usage/help for all Makefile commands
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## confirm: Ask for confirmation before continuing
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ "$$ans" = "y" ]


# ========================================================================================== #
# QUALITY CONTROL
# ========================================================================================== #

.PHONY: audit vendor

## audit: tidy dependencies and format , vet and test all code
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
vendor:
	@echo 'Tidying and verifying module  dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# =================================================================================== #
# DEVELOPMENT
# =================================================================================== #

.PHONY: run/api db/migrations/new db/migrations/up db/migrations/goto db/migrations/down db/migrations/rollback db/migrations/force docs/gen test

## run/api: Run the main Go API server
run/api:
	go run ./cmd/api -db-dsn=${DB_DSN}

## db/migrations/new: Create a new migration file (provide name with 'name=...')
db/migrations/new:
	migrate create -seq -ext sql -dir ${MIGRATIONS_PATH} $(name)

## db/migrations/up: Apply all up migrations to the database
db/migrations/up: confirm
	migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" up

## db/migrations/goto: Migrate the database to a specific version (provide with 'version=...')
db/migrations/goto: confirm
	migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" goto $(version)

## db/migrations/down: Apply all down migrations (rollback all)
db/migrations/down: confirm
	migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" down

## db/migrations/rollback: Rollback a specific number of migrations (provide with 'version=...')
db/migrations/rollback: confirm
	migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" down $(version)

## db/migrations/force: Force the database schema to a specific version (provide with 'version=...')
db/migrations/force: confirm
	migrate -path=${MIGRATIONS_PATH} -database="${DB_DSN}" force $(version)

## docs/gen: Generate API documentation using swag
docs/gen:
	@echo "Generating API documentation..."
	swag init -g ./api/main.go -d cmd,internal && swag fmt
	@echo "API documentation generated successfully."



# =================================================================================== #
# BUILD
# =================================================================================== #

.PHONY: build/api


## build/api: build the cmp/api application
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
