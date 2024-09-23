.PHONY: all \
	serve \
	ti test-integration \
	swag swagger \
	sch schema \
	gen generate \
	mig migrate \
	seeds \
	db-up database-up \
	db-down database-down

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

seeds:
	go run scripts/seeds/dw.go

db-up: database-up
database-up:
	docker-compose up -d

db-down: database-down
database-down:
	docker-compose down
