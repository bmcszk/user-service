include .env

.PHONY: install-migrate
install-migrate:
	(migrate -version) || (go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)

.PHONY: migrate-up
migrate-up: install-migrate
	migrate -path db/migrations -database "${POSTGRES_URL}" up

.PHONY: migrate-down
migrate-down: install-migrate
	migrate -path db/migrations -database "${POSTGRES_URL}" down 1

.PHONY: install-sqlc
install-sqlc:
	(sqlc version) || (go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)

.PHONY: sqlc
sqlc: install-sqlc
	sqlc generate -f db/sqlc.yaml

.PHONY: compose
compose:
	docker compose up -d --build --wait

.PHONY: test
test: compose
	go test -count=1 ./...

.PHONY: coverage
coverage: compose
	go test -count=1 -coverprofile coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
# TODO figure out how to do coveage during e2e tests

include tilt/Makefile
