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

s: serve
serve:
	air

t: test
test:
	go test -v $$(go list ./... | grep -v 'ent/\|docs/\|_integration_test.go')

ti: test-integration
test-integration:
	$(MAKE) gen
	go test -p 1 -v $$(go list ./... | grep -v 'ent/\|docs/') -tags=production

swag: swagger
swagger:
	swag init

sch:
	go run scripts/schema/main.go $(filter-out $@,$(MAKECMDGOALS))
	@:
schema:
	go run scripts/schema/main.go $(filter-out $@,$(MAKECMDGOALS))
	@:

gen: generate
generate:
	go run scripts/generate/main.go

mig: migrate
migrate:
	go run scripts/migrate/main.go

seeds:
	go run scripts/seeds/main.go $(filter-out $@,$(MAKECMDGOALS))
	@:

sy:
	go run scripts/seeds/main.go $(filter-out $@,$(MAKECMDGOALS)) -y
	@:
seeds-y:
	go run scripts/seeds/main.go $(filter-out $@,$(MAKECMDGOALS)) -y
	@:

db-up: database-up
database-up:
	docker-compose up -d

db-down: database-down
database-down:
	docker-compose down

h: hooks
hooks:
	husky install
