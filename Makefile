# ========================================================================= #
# HELPERS
# ========================================================================= #

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ========================================================================= #
# DEVELOPMENT
# ========================================================================= #

## run/api: run the cmd/api application
run/api:
	# go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}
	go run ./cmd/api

## db/psql: connect to the database using psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration 
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -ext=.sql -dir=./migrations -seq ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	# migrate -path=./migrations -database "postgres://greenlight:greenlight@localhost/greenlight?sslmode=disable" -verbose up
	migrate -path=./migrations -database ${GREENLIGHT_DB_DSN} -verbose up

## db/psql: apply all down database migrations
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path=./migrations -database "postgres://greenlight:greenlight@localhost/greenlight?sslmode=disable" -verbose down

# ========================================================================= #
# QUALITY CONTROL
# ========================================================================= #
## audit: tidy dependencies and	format code, vet and test all code
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

.PHONY : audit help vendor confirm run/api db/psql db/migrate/up db/migrate/down db/migration