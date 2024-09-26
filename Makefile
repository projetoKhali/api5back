all: serve

%:
	@:

serve:
	air

t: test
test:
	go test -v $$(go list ./... | grep -v 'ent/\|docs/\|_integration_test.go')

ti: test-integration
test-integration:
	go test -v $$(go list ./... | grep -v 'ent/\|docs/') -tags=integration

swag: swagger
swagger:
	swag init

sch: schema
schema:
	go run scripts/schema/main.go $(filter-out $@,$(MAKECMDGOALS))

gen: generate
generate:
	go run scripts/generate/main.go

mig: migrate
migrate:
	go run scripts/migrate/main.go

db-up: database-up
database-up:
	docker-compose up -d

db-down: database-down
database-down:
	docker-compose down

h: hooks
hooks:
	husky install
