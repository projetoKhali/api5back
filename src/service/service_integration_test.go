//go:build integration
// +build integration

package service

import (
	"context"
	"testing"

	"api5back/ent"
	"api5back/ent/dimvacancy"
	"api5back/seeds"
	"api5back/src/database"
	"api5back/src/property"

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
		var testVacancyDbId int
		var hiringProcessCandidateId int

		for _, TestCase := range []database.TestCase{
			{
				Name: "Insert a hiring_process_candidate into the table",
				Run: func(t *testing.T) {
					dimVacancy, err := intEnv.
						Client.
						DimVacancy.
						Query().
						First(ctx)
					require.NoError(t, err)

					hiringProcessCandidate, err := intEnv.Client.HiringProcessCandidate.
						Create().
						SetDimVacancyDbId(dimVacancy.DbId).
						SetName("John Doe").
						SetEmail("John@Doe.com").
						SetPhone("+1234567890").
						SetApplyDate(dimVacancy.OpeningDate).
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
					dimVacancy, err := intEnv.
						Client.
						DimVacancy.
						Query().
						WithHiringProcessCandidates().
						Where(dimvacancy.DbId(testVacancyDbId)).
						Order(
							ent.Desc(dimvacancy.FieldClosingDate),
							ent.Desc(dimvacancy.FieldID),
						).
						First(ctx)
					require.NoError(t, err)

					require.NoError(t, err)

					candidates, err := dimVacancy.
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

	if testResult := t.Run("GetMetrics returns correct metrics", func(t *testing.T) {
		metricsData, err := GetMetrics(
			ctx, intEnv.Client,
			FactHiringProcessFilter{
				Recruiters:    []int{},
				Processes:     []int{},
				Vacancies:     []int{},
				DateRange:     nil,
				ProcessStatus: []int{},
				VacancyStatus: []int{},
			},
		)

		require.NoError(t, err)

		require.NotNil(t, metricsData)

		require.NotEmpty(t, metricsData.CardInfos)
		require.NotNil(t, metricsData.VacancySummary)
		require.NotNil(t, metricsData.AverageHiringTime)
	}); !testResult {
		t.Fatalf("GetMetrics test failed")
	}
}

func TestTableDashboard(t *testing.T) {
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

	if testResult := t.Run("Vacancy Table returns all FactHiringProcess", func(t *testing.T) {
		vacancies, err := GetVacancyTable(
			ctx, intEnv.Client,
			FactHiringProcessFilter{
				Recruiters:    []int{},
				Processes:     []int{},
				Vacancies:     []int{},
				DateRange:     nil,
				ProcessStatus: []int{},
				VacancyStatus: []int{},
				Page:          nil,
				PageSize:      nil,
			},
		)

		require.NoError(t, err)
		require.NotNil(t, vacancies)
		require.Equal(t, 5, len(vacancies.FactHiringProcess))
	}); !testResult {
		t.Fatalf("GetVacancyTable no filter test failed")
	}

	if testResult := t.Run("Vacancy Table returns correct number of FactHiringProcess", func(t *testing.T) {
		vacancies, err := GetVacancyTable(
			ctx, intEnv.Client,
			FactHiringProcessFilter{
				Recruiters: []int{},
				Processes:  []int{},
				Vacancies:  []int{},
				DateRange: &DateRange{
					StartDate: "2024-07-16",
					EndDate:   "2024-08-12",
				},
				ProcessStatus: []int{},
				VacancyStatus: []int{},
				Page:          nil,
				PageSize:      nil,
			},
		)

		require.NoError(t, err)
		require.NotNil(t, vacancies)
		require.Equal(t, 1, len(vacancies.FactHiringProcess))
	}); !testResult {
		t.Fatalf("GetVacancyTable dateRange test failed")
	}
}
