set dotenv-load

default:
    @just --list

run_dev:
    @go run .

migration_up:
	@migrate -path ./migrations -database "$IGN_DB_URL_NOSSL" up

migration_down:
	@migrate -path ./migrations -database "$IGN_DB_URL_NOSSL" down
