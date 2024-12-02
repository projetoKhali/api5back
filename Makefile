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

# Test with optional module parameter
t: test
test:
	@if [ -z "$(module)" ]; then \
		go test -v $$(go list ./... | grep -v 'ent/\|docs/\|_integration_test.go'); \
	else \
		go test -v api5back/src/$(module); \
	fi

# Test integration with optional module parameter
ti: test-integration
test-integration:
	$(MAKE) gen
	@if [ -z "$(module)" ]; then \
		go test -p 1 -v $$(go list ./... | grep -v 'ent/\|docs/') -tags=integration; \
	else \
		go test -p 1 -v api5back/src/$(module) -tags=integration; \
	fi

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
