//go:build integration
// +build integration

package service

import (
	"context"
	"fmt"
	"testing"

	"api5back/seeds"
	"api5back/src/database"
	"api5back/src/model"

	"github.com/stretchr/testify/require"
)

func TestGetSuggestionsFunctions(t *testing.T) {
	ctx := context.Background()
	var intEnv *database.IntegrationEnvironment = nil

	if testResult := t.Run("Setup database connection", func(t *testing.T) {
		intEnv = database.
			DefaultIntegrationEnvironment(ctx).
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
		ExpectedLength     int
		ExpectedError      error
	}

	pageSize := 1000
	maxPageSizeRequest := model.SuggestionsPageRequest{
		PageRequest: &model.PageRequest{
			Page:     nil,
			PageSize: &pageSize,
		},
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
			ExpectedLength: len(seeds.DwDimUser),
			ExpectedError:  nil,
		},
		{
			Name: "GetProcessSuggestions returns all unique processes by DbId",
			GetSuggestionsFunc: func() (int, error) {
				suggestions, err := GetProcessSuggestions(
					ctx,
					intEnv.Client,
					&model.SuggestionsFilter{
						SuggestionsPageRequest: maxPageSizeRequest,
					},
				)
				return len(suggestions.Items), err
			},
			ExpectedLength: len(seeds.DwDimProcess),
			ExpectedError:  nil,
		},
		{
			Name: "GetVacancySuggestions returns all unique vacancies by DbId",
			GetSuggestionsFunc: func() (int, error) {
				fmt.Printf("maxPageSizeRequest: %+v | pageSize: %+v\n", maxPageSizeRequest, pageSize)
				suggestions, err := GetVacancySuggestions(
					ctx,
					intEnv.Client,
					&model.SuggestionsFilter{
						SuggestionsPageRequest: maxPageSizeRequest,
					},
				)
				return len(suggestions.Items), err
			},
			ExpectedLength: len(seeds.DwDimVacancy),
			ExpectedError:  nil,
		},
	} {
		if testResult := t.Run(testCase.Name, func(t *testing.T) {
			suggestionsCount, err := testCase.GetSuggestionsFunc()

			require.Equal(t, testCase.ExpectedError, err)

			if testCase.ExpectedError == nil {
				require.Equal(t, testCase.ExpectedLength, suggestionsCount)
			}
		}); !testResult {
			t.Fatalf("Test case failed")
		}
	}
}
