all: serve

%:
	@:

serve:
	air

schema:
	go run scripts/schema/main.go $(filter-out $@,$(MAKECMDGOALS))

gen:
	go run scripts/generate/main.go

migrate:
	go run scripts/migrate/main.go

db-up:
	docker-compose up -d

db-down:
	docker-compose down
