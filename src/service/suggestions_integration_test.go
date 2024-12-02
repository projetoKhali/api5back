//go:build integration
// +build integration

package service

import (
	"context"
	"math"
	"testing"

	"api5back/ent"
	"api5back/seeds"
	"api5back/src/database"
	"api5back/src/model"

	"github.com/stretchr/testify/require"
)

func TestGetSuggestionsFunctions(t *testing.T) {
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

	type TestFunc func() (int, error)
	type TestCase struct {
		Name               string
		GetSuggestionsFunc TestFunc
		ExpectedFunc       TestFunc
	}

	maxPageSizeRequest := model.PageRequest{
		Page:     nil,
		PageSize: &[]int{math.MaxInt32}[0],
	}

	for _, testCase := range []TestCase{
		{
			Name: "GetUserSuggestions returns all unique users by DbId",
			GetSuggestionsFunc: func() (int, error) {
				suggestions, err := GetUserSuggestions(
					ctx,
					intEnv.Client,
					&maxPageSizeRequest,
				)
				return len(suggestions.Items), err
			},
			ExpectedFunc: func() (int, error) {
				wrappedUsers := make([]DbIdGetter, len(seeds.DwDimUser))

				for i, user := range seeds.DwDimUser {
					wrappedUsers[i] = &DimUserWrapper{&user}
				}

				return len(deduplicateById(wrappedUsers)), nil
			},
		},
		{
			Name: "GetProcessSuggestions returns all unique processes by DbId",
			GetSuggestionsFunc: func() (int, error) {
				suggestions, err := GetProcessSuggestions(
					ctx,
					intEnv.Client,
					&model.SuggestionsPageRequest{
						PageRequest: &maxPageSizeRequest,
					},
				)
				return len(suggestions.Items), err
			},
			ExpectedFunc: func() (int, error) {
				wrappedProcesses := make([]DbIdGetter, len(seeds.DwDimProcess))

				for i, process := range seeds.DwDimProcess {
					wrappedProcesses[i] = &DimProcessWrapper{&process}
				}

				return len(deduplicateById(wrappedProcesses)), nil
			},
		},
		{
			Name: "GetVacancySuggestions returns all unique vacancies by DbId",
			GetSuggestionsFunc: func() (int, error) {
				suggestions, err := GetVacancySuggestions(
					ctx,
					intEnv.Client,
					&model.SuggestionsPageRequest{
						PageRequest: &maxPageSizeRequest,
					},
				)
				return len(suggestions.Items), err
			},
			ExpectedFunc: func() (int, error) {
				wrappedVacancies := make([]DbIdGetter, len(seeds.DwDimVacancy))

				for i, vacancy := range seeds.DwDimVacancy {
					wrappedVacancies[i] = &DimVacancyWrapper{&vacancy}
				}

				return len(deduplicateById(wrappedVacancies)), nil
			},
		},
	} {
		if testResult := t.Run(testCase.Name, func(t *testing.T) {
			suggestionsCount, err := testCase.GetSuggestionsFunc()
			expectedCount, expectedErr := testCase.ExpectedFunc()

			require.Equal(t, expectedErr, err)
			require.Equal(t, expectedCount, suggestionsCount)
		}); !testResult {
			t.Fatalf("Test case failed")
		}
	}
}

type DbIdGetter interface{ GetDbId() int }

type (
	DimUserWrapper    struct{ *ent.DimUser }
	DimVacancyWrapper struct{ *ent.DimVacancy }
	DimProcessWrapper struct{ *ent.DimProcess }
)

func (duw *DimUserWrapper) GetDbId() int    { return duw.DbId }
func (dwv *DimVacancyWrapper) GetDbId() int { return dwv.DbId }
func (dpw *DimProcessWrapper) GetDbId() int { return dpw.DbId }

func deduplicateById(suggestions []DbIdGetter) []DbIdGetter {
	seen := make(map[int]bool)
	uniqueSuggestions := []DbIdGetter{}

	for _, suggestion := range suggestions {
		id := suggestion.GetDbId()

		if !seen[id] {
			seen[id] = true
			uniqueSuggestions = append(uniqueSuggestions, suggestion)
		}
	}

	return uniqueSuggestions
}
