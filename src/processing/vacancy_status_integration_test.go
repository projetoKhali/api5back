//go:build integration
// +build integration

package processing

import (
	"context"
	"testing"

	"api5back/ent"
	"api5back/seeds"
	"api5back/src/database"

	"github.com/stretchr/testify/require"
)

func TestVacancyStatusProcessing(t *testing.T) {
	ctx := context.Background()
	var intEnv *database.IntegrationEnvironment
	var err error

	if testResult := t.Run("Setup database connection", func(t *testing.T) {
		intEnv = database.DefaultIntegrationEnvironment(ctx).
			WithSeeds(seeds.DataWarehouse)

		require.NotNil(t, intEnv)
		require.NoError(t, intEnv.Error)
		require.NotNil(t, intEnv.Client)
	}); !testResult {
		t.Fatalf("Setup test failed")
	}

	var factHiringProcesses []*ent.FactHiringProcess

	if testResult := t.Run("query FactHiringProcess", func(t *testing.T) {
		factHiringProcesses, err = intEnv.
			Client.
			FactHiringProcess.
			Query().
			WithDimDatetime().
			WithDimProcess().
			WithDimUser().
			WithDimVacancy().
			All(ctx)

		require.NoError(t, err)
	}); !testResult {
		t.Fatalf("failed to query FactHiringProcess: %v", err)
	}

	if testResult := t.Run("Test vacancy status processing", func(t *testing.T) {
		vacancyStatusSummary, err := GenerateVacancyStatusSummary(factHiringProcesses)
		require.NoError(t, err)

		require.Equal(t,
			VacancyStatusSummary{
				Open:      9,
				Analyzing: 0,
				Closed:    1,
			},
			vacancyStatusSummary,
		)
	}); !testResult {
		t.Fatalf("Vacancy status processing test failed")
	}
}
