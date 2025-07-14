set windows-powershell := true
set dotenv-load
set dotenv-required

app_dir := "./cmd/server"

default:
    @just --list

build-frontend:
    @templ generate
    @npm run build

[private]
build-air: build-frontend
    @go build -o ./tmp_air/main.exe "{{ app_dir }}/main.go"

build-dev: build-frontend
    @go build "{{ app_dir }}"

build-prod: build-frontend
    @go build -tags prod "{{ app_dir }}"

run: build-frontend
    @go run "{{ app_dir }}/main.go"

dev:
    @air -c scripts/air.toml

setup-dev-env:
    @go install github.com/air-verse/air@v1.61.5
    @go install github.com/a-h/templ/cmd/templ@v0.3.856
    @go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
    @go mod tidy
    @npm i

build-docker tag:
    @ docker build -t {{tag}} -f docker/Dockerfile .

compose-up:
    @docker compose -f docker/docker-compose.yaml up

compose-down:
    @docker compose -f docker/docker-compose.yaml down

db-migrate:
    @goose pgx -dir ./schema/migrations up

db-seed:
    @goose pgx -dir ./schema/seed -no-versioning up