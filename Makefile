all: serve

serve:
	air

ent-gen:
	go run scripts/generate/main.go

ent-migrate:
	go run scripts/migrate/main.go

db-up:
	docker-compose up -d

db-down:
	docker-compose down
