//go:build integration
// +build integration

package service

import (
	"context"
	"testing"

	"api5back/ent/facthiringprocess"
	"api5back/seeds"
	"api5back/src/database"
	"api5back/src/property"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestDatabaseOperations(t *testing.T) {
	ctx := context.Background()
	var intEnv *database.IntegrationEnvironment = nil

	if testResult := t.Run("Setup database connection", func(t *testing.T) {
		intEnv = database.DefaultIntegrationEnvironment(ctx).
			WithSeeds(seeds.DataWarehouse)

		require.NotNil(t, intEnv)
		require.NoError(t, intEnv.Error)
		require.NotNil(t, intEnv.Client)
	}); !testResult {
		t.Fatalf("Setup test failed")
	}

	t.Run("Test hiring_process_candidate table operations", func(t *testing.T) {
		var testFactHiringProcessId int
		var hiringProcessCandidateId int

		for _, TestCase := range []database.TestCase{
			{
				Name: "Insert a hiring_process_candidate into the table",
				Run: func(t *testing.T) {
					factHiringProcess, err := intEnv.
						Client.
						FactHiringProcess.
						Query().
						WithDimVacancy().
						First(ctx)
					require.NoError(t, err)

					testFactHiringProcessId = factHiringProcess.ID

					dimVacancy, err := factHiringProcess.
						Edges.
						DimVacancyOrErr()
					require.NoError(t, err)

					applyDatePgType := &pgtype.Date{}
					if err := applyDatePgType.Scan(dimVacancy.OpeningDate); err != nil {
						t.Fatalf("failed to parse applyDatePgType for candidate: %v", err)
					}

					hiringProcessCandidate, err := intEnv.Client.HiringProcessCandidate.
						Create().
						SetFactHiringProcessID(testFactHiringProcessId).
						SetName("John Doe").
						SetEmail("John@Doe.com").
						SetPhone("+1234567890").
						SetApplyDate(applyDatePgType).
						SetStatus(property.HiringProcessCandidateStatusInAnalysis).
						SetScore(0).
						Save(ctx)
					if err != nil {
						t.Fatalf("failed to insert the hiring_process_candidate: %v", err)
					}

					hiringProcessCandidateId = hiringProcessCandidate.ID
				},
			},
			{
				Name: "Select candidate list from the edges of a FactHiringProcess",
				Run: func(t *testing.T) {
					factHiringProcesses, err := intEnv.
						Client.
						FactHiringProcess.
						Query().
						WithDimVacancy().
						WithHiringProcessCandidates().
						Where(facthiringprocess.ID(testFactHiringProcessId)).
						First(ctx)
					require.NoError(t, err)

					candidates, err := factHiringProcesses.
						Edges.
						HiringProcessCandidatesOrErr()
					require.NoError(t, err)
					require.NotNil(t, candidates)
					require.NotEmpty(t, candidates)
				},
			},
			{
				Name: "Select a candidate by ID",
				Run: func(t *testing.T) {
					hiringProcessCandidate, err := intEnv.
						Client.
						HiringProcessCandidate.
						Get(ctx, hiringProcessCandidateId)
					require.NoError(t, err)
					require.NotNil(t, hiringProcessCandidate)
					require.Equal(
						t,
						property.HiringProcessCandidateStatusInAnalysis,
						hiringProcessCandidate.Status,
					)
				},
			},
		} {
			if testResult := t.Run(TestCase.Name, TestCase.Run); !testResult {
				t.Fatalf("Test case failed")
			}
		}
	})

	metricsService := NewMetricsService(intEnv.Client)
	t.Run("GetMetrics returns correct metrics", func(t *testing.T) {
		metricsData, err := metricsService.GetMetrics(ctx)

		require.NoError(t, err)

		require.NotNil(t, metricsData)

		require.NotEmpty(t, metricsData.CardInfos)
		require.NotNil(t, metricsData.VacancySummary)
		require.NotNil(t, metricsData.AvgHiringTime)

	})
}
