package database

import (
	"database/sql"
	"fmt"

	"api5back/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Setup(prefix string) (*ent.Client, error) {
	databaseCredentials, err := newDatabaseCredentials(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create databaseCredentials: %v", err)
	}
	databaseUrl := databaseCredentials.getConnectionString()

	client, err := createPostgresClient(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %v", err)
	}

	return client, nil
}

func createPostgresClient(databaseUrl string) (*ent.Client, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}
