//go:build integration
// +build integration

package database

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"api5back/ent"
	"api5back/ent/migrate"

	"github.com/stretchr/testify/require"
)

func TestDatabaseOperations(t *testing.T) {
	ctx := context.Background()
	var intEnv *IntegrationEnvironment = nil
	var err error

	if testResult := t.Run("Setup database connection", func(t *testing.T) {
		intEnv = DefaultIntegrationEnvironment(ctx)
		require.NotNil(t, intEnv)
		require.NoError(t, intEnv.Error)
		require.NotNil(t, intEnv.Client)
		println("a")
	}); !testResult {
		t.Fatalf("Setup test failed")
	}

	if testResult := t.Run("Migrate database", func(t *testing.T) {
		if err = intEnv.Client.Schema.Create(
			ctx,
			migrate.WithDropIndex(true),
			migrate.WithDropColumn(true),
		); err != nil {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("failed to migrate the database: %v", err))
			sb.WriteString("\n\nThis error may be caused by the test not waiting long enough for the database to be ready.")
			sb.WriteString("\nTry increasing the sleep time in the `.env.integration` test.")
			t.Fatalf(sb.String())
		}
	}); !testResult {
		t.Fatalf("Migration test failed")
	}

	t.Run("Test dim_user table operations", func(t *testing.T) {
		var testDimUser *ent.DimUser

		for _, TestCase := range []TestCase{
			{
				Name: "Insert a dim_user into the table",
				Run: func(t *testing.T) {
					testDimUser, err = intEnv.Client.DimUser.
						Create().
						SetName("John Doe").
						SetOcupation("Software Engineer").
						Save(ctx)
					if err != nil {
						t.Fatalf("failed to insert the dim_user: %v", err)
					}
					require.Equal(t, "John Doe", testDimUser.Name)
					require.Equal(t, "Software Engineer", testDimUser.Ocupation)
				},
			}, {
				Name: "Retrieve the inserted dim_user",
				Run: func(t *testing.T) {
					retrievedDimUser, err := intEnv.Client.DimUser.Get(ctx, testDimUser.ID)
					if err != nil {
						t.Fatalf("failed to retrieve the dim_user: %v", err)
					}
					require.Equal(t, testDimUser.ID, retrievedDimUser.ID)
					require.Equal(t, testDimUser.Name, retrievedDimUser.Name)
					require.Equal(t, testDimUser.Ocupation, retrievedDimUser.Ocupation)
				},
			}, {
				Name: "Delete the dim_user",
				Run: func(t *testing.T) {
					err = intEnv.Client.DimUser.DeleteOne(testDimUser).Exec(ctx)
					require.NoError(t, err)
				},
			}, {
				Name: "Try to retrieve the dim_user again, expecting a not found error",
				Run: func(t *testing.T) {
					_, err = intEnv.Client.DimUser.Get(ctx, testDimUser.ID)
					require.Error(t, err)
				},
			},
		} {
			if testResult := t.Run(TestCase.Name, TestCase.Run); !testResult {
				t.Fatalf("Test case failed")
			}
		}
	})
}
