package main

import (
	"api5back/ent"
	"api5back/src/database"
	"context"
	"fmt"
)

var databasePrefixes = []string{
	"DB",
	"DW",
}

// Run the automatic migration tool to create all schema resources of the database.
func Migrate(client *ent.Client) error {
	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("scripts/migrate • failed creating schema resources: %v", err)
	}

	return nil
}

// Run the migration tool of all databases.
func MigrateAll() error {
	for _, prefix := range databasePrefixes {
		client, err := database.Setup(prefix)
		if err != nil {
			return fmt.Errorf("scripts/migrate • failed to setup database: %v", err)
		}
		defer client.Close()

		if err := Migrate(client); err != nil {
			return fmt.Errorf("scripts/migrate • failed to migrate database: %v", err)
		}

		fmt.Printf("scripts/migrate • migrated database with prefix: %s\n", prefix)
	}

	return nil
}

// manual entry point for migration on command
func main() {
	fmt.Println("scripts/migrate • Migrating all databases...")

	if err := MigrateAll(); err != nil {
		panic(fmt.Errorf("scripts/migrate • failed to migrate databases: %v", err))
	}

	fmt.Println("scripts/migrate • Successfully migrated databases.")
}