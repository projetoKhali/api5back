//go:build integration
// +build integration

package database

import (
	"api5back/ent"
	"api5back/ent/migrate"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Generate the ent files
	err := entc.Generate("../schema", &gen.Config{
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

	client := ent.NewClient(
		ent.Driver(
			entsql.OpenDB(dialect.Postgres, db),
		),
	)

	// Migrate the database
	ctx := context.Background()
	err = client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	require.NoError(t, err)

	// Insert a user into the table
	user, err := client.DimUser.
		Create().
		SetName("John Doe").
		SetOcupation("Software Engineer").
		Save(ctx)
	require.NoError(t, err)
	require.Equal(t, "John Doe", user.Name)
	require.Equal(t, "Software Engineer", user.Ocupation)

	// Retrieve the inserted user
	retrievedUser, err := client.DimUser.Get(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ID, retrievedUser.ID)
	require.Equal(t, user.Name, retrievedUser.Name)
	require.Equal(t, user.Ocupation, retrievedUser.Ocupation)

	// Delete the user
	err = client.DimUser.DeleteOne(user).Exec(ctx)
	require.NoError(t, err)

	// Try to retrieve the user again, expecting a not found error
	_, err = client.DimUser.Get(ctx, user.ID)
	require.Error(t, err)
}
