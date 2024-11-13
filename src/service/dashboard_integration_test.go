//go:build integration
// +build integration

package service

import (
	"context"
	"testing"

	"api5back/seeds"
	"api5back/src/database"
	"api5back/src/model"

	"github.com/stretchr/testify/require"
)

func TestDashboard(t *testing.T) {
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

	if testResult := t.Run("GetMetrics returns correct metrics", func(t *testing.T) {
		metricsData, err := GetMetrics(
			ctx, intEnv.Client,
			model.FactHiringProcessFilter{
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
			model.FactHiringProcessFilter{
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
		require.Equal(t, 10, len(vacancies.Items))
	}); !testResult {
		t.Fatalf("GetVacancyTable no filter test failed")
	}

	if testResult := t.Run("Vacancy Table returns correct number of FactHiringProcess", func(t *testing.T) {
		vacancies, err := GetVacancyTable(
			ctx, intEnv.Client,
			model.FactHiringProcessFilter{
				Recruiters: []int{},
				Processes:  []int{},
				Vacancies:  []int{},
				DateRange: &model.DateRange{
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
		require.Equal(t, 1, len(vacancies.Items))
	}); !testResult {
		t.Fatalf("GetVacancyTable dateRange test failed")
	}
}
