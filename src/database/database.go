package database

import (
	"api5back/ent"
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func openConnection(databaseUrl string) (*ent.Client, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}

func DatabaseSetup() error {
	// Run the ent codegen tool to generate the schema files.
	if err := entc.Generate("./ent/schema", &gen.Config{}); err != nil {
		return fmt.Errorf("error running ent codegen: %w", err)
	}

	// Retrieve the database credentials from environment variables.
	databaseCredentials, err := NewDatabaseCredentials()
	if err != nil {
		return fmt.Errorf("failed to create databaseCredentials: %v", err)
	}
	connectionString := databaseCredentials.GetConnectionString()

	client, err := openConnection(connectionString)

	ctx := context.Background()

	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
	}

	return nil
}
