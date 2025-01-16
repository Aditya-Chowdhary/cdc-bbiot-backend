#!make
include .env

# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	migrate create -seq -ext=.sql -dir=./internal/database/migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	migrate -path=./internal/database/migrations -database="${DB_URL}?sslmode=disable" up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	migrate -path=./internal/database/migrations -database="${DB_URL}?sslmode=disable" down

## migrations/goto version=$1: migrate to a specific version number
.PHONY: migrations/goto
migrations/goto:
	migrate -path=./internal/database/migrations -database="${DB_URL}?sslmode=disable" goto ${version}

## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	migrate -path=./internal/database/migrations -database="${DB_URL}?sslmode=disable" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	migrate -path=./internal/database/migrations -database="${DB_URL}?sslmode=disable" version
