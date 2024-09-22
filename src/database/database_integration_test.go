//go:build integration
// +build integration

package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"api5back/ent"
	"api5back/ent/migrate"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestDatabaseOperations(t *testing.T) {
	ctx := context.Background()
	var client *ent.Client
	var err error

	t.Run("Setup database connection", func(t *testing.T) {
		// Generate the ent files
		err = entc.Generate("../schema", &gen.Config{
			Schema:  "../schema",
			Target:  "../../ent",
			Package: "api5back/ent",
		})
		require.NoError(t, err)

		// Load the environment variables
		err = godotenv.Load("../../.env")
		require.NoError(t, err)

		// Create the database client
		databaseUrl := fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			os.Getenv("DW_USER"),
			os.Getenv("DW_PASS"),
			os.Getenv("DW_HOST"),
			os.Getenv("DW_PORT"),
			os.Getenv("DW_NAME"),
		)

		db, err := sql.Open("pgx", databaseUrl)
		require.NoError(t, err)

		client = ent.NewClient(
			ent.Driver(
				entsql.OpenDB(dialect.Postgres, db),
			),
		)
	})

	defer client.Close()

	t.Run("Migrate database", func(t *testing.T) {
		err = client.Schema.Create(
			ctx,
			migrate.WithDropIndex(true),
			migrate.WithDropColumn(true),
		)
		require.NoError(t, err)
	})

	t.Run("Test dim_user table operations", func(t *testing.T) {
		var testDimUser *ent.DimUser

		t.Run("Insert a dim_user into the table", func(t *testing.T) {
			testDimUser, err = client.DimUser.
				Create().
				SetName("John Doe").
				SetOcupation("Software Engineer").
				Save(ctx)
			require.NoError(t, err)
			require.Equal(t, "John Doe", testDimUser.Name)
			require.Equal(t, "Software Engineer", testDimUser.Ocupation)
		})

		t.Run("Retrieve the inserted dim_user", func(t *testing.T) {
			retrievedDimUser, err := client.DimUser.Get(ctx, testDimUser.ID)
			require.NoError(t, err)
			require.Equal(t, testDimUser.ID, retrievedDimUser.ID)
			require.Equal(t, testDimUser.Name, retrievedDimUser.Name)
			require.Equal(t, testDimUser.Ocupation, retrievedDimUser.Ocupation)
		})

		t.Run("Delete the dim_user", func(t *testing.T) {
			err = client.DimUser.DeleteOne(testDimUser).Exec(ctx)
			require.NoError(t, err)
		})

		t.Run("Try to retrieve the dim_user again, expecting a not found error", func(t *testing.T) {
			_, err = client.DimUser.Get(ctx, testDimUser.ID)
			require.Error(t, err)
		})
	})
}
