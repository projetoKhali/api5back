all: serve

%:
	@:

serve:
	air

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
