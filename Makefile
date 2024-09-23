all: serve

%:
	@:

serve:
	air

ti: test-integration
test-integration:
	$(MAKE) gen
	go test -p 1 -v $$(go list ./... | grep -v 'ent/\|docs/') -tags=integration

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
