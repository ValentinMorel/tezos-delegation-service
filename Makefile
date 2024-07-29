.PHONY: all


show-doc:
	open docs/index.html

show-coverage: 
	open coverage.html

up:
	PGPASSWORD="password" psql -h localhost tezos -U username -p 5432 -f db/migrations/000001_delegations.up.sql

down:
	PGPASSWORD="password" psql -h localhost tezos -U username -p 5432 -f db/migrations/000001_delegations.down.sql


test-db:
	PGPASSWORD="password" psql -h localhost tezos -U username -p 5432  -c "CREATE DATABASE db_test" || exit 0

test-up:
	PGPASSWORD="password" psql -h localhost db_test -U username -p 5432 -f db/migrations/000001_delegations.up.sql

test-down:
	PGPASSWORD="password" psql -h localhost db_test -U username -p 5432 -f db/migrations/000001_delegations.down.sql

run: 
	go run cmd/server/main.go

generate: 
	oapi-codegen --generate gin,spec --package api -o internal/api/openapi.gen.go openapi.yaml &&\
	cd db && \
	sqlc generate
	