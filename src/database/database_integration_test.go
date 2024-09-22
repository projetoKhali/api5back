//go:build integration
// +build integration

package database

import (
	"context"
	"testing"

	"api5back/ent"
	"api5back/ent/migrate"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
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

		databaseCredentials, err := newTestingDatabaseCredentials()
		require.NoError(t, err)

		databaseUrl := databaseCredentials.getConnectionString()

		client, err = createPostgresClient(databaseUrl)
		require.NoError(t, err)
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
