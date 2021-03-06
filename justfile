set dotenv-load

default:
    @just --list

run:
    @go run . --dev=prod

run_dev:
    @go run .

migration_up:
	@migrate -path ./migrations -database "$DATABASE_URL_NOSSL" up

migration_down:
	@migrate -path ./migrations -database "$DATABASE_URL_NOSSL" down

populate_fake_data:
    @psql -U "$DATABASE_USER" -d "$DATABASE_NAME" -f ./resources/fakedata.sql
