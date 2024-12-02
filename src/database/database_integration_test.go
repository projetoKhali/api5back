//go:build integration
// +build integration

package database

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"api5back/ent"
	"api5back/ent/facthiringprocess"
	"api5back/ent/migrate"
	"api5back/src/property"

	"github.com/stretchr/testify/require"
)

func TestBaseDatabaseOperations(t *testing.T) {
	ctx := context.Background()
	var intEnv *IntegrationEnvironment = nil
	var err error

	if testResult := t.Run("Setup database connection", func(t *testing.T) {
		intEnv = DefaultIntegrationEnvironment(ctx)

		require.NotNil(t, intEnv)
		require.NoError(t, intEnv.Error)
		require.NotNil(t, intEnv.Client)
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
						SetDbId(1).
						SetName("John Doe").
						SetOccupation("Software Engineer").
						Save(ctx)
					if err != nil {
						t.Fatalf("failed to insert the dim_user: %v", err)
					}
					require.Equal(t, "John Doe", testDimUser.Name)
					require.Equal(t, "Software Engineer", testDimUser.Occupation)
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
					require.Equal(t, testDimUser.Occupation, retrievedDimUser.Occupation)
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

	t.Run("Test hiring_process_candidate table operations", func(t *testing.T) {
		var testFactHiringProcessId int
		var hiringProcessCandidateId int

		for _, TestCase := range []TestCase{
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

					dimVacancy, err := factHiringProcess.
						Edges.
						DimVacancyOrErr()
					require.NoError(t, err)

					dimCandidate, err := intEnv.
						Client.
						DimCandidate.
						Create().
						SetDimVacancyDbId(factHiringProcess.DimVacancyId).
						SetDbId(1).
						SetName("John Doe").
						SetEmail("John@Doe.com").
						SetPhone("+1234567890").
						SetApplyDate(dimVacancy.OpeningDate).
						SetStatus(property.DimCandidateStatusInAnalysis).
						SetScore(0).
						Save(ctx)
					if err != nil {
						t.Fatalf("failed to insert the hiring_process_candidate: %v", err)
					}

					hiringProcessCandidateId = dimCandidate.ID
				},
			},
			{
				Name: "Select candidate list from the edges of a DimVacancy",
				Run: func(t *testing.T) {
					factHiringProcesses, err := intEnv.
						Client.
						FactHiringProcess.
						Query().
						WithDimVacancy().
						Where(facthiringprocess.ID(testFactHiringProcessId)).
						First(ctx)
					require.NoError(t, err)

					dimVacancy, err := factHiringProcesses.
						Edges.
						DimVacancyOrErr()
					require.NoError(t, err)

					candidates, err := dimVacancy.
						Edges.
						DimCandidatesOrErr()
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
						DimCandidate.
						Get(ctx, hiringProcessCandidateId)
					require.NoError(t, err)
					require.NotNil(t, hiringProcessCandidate)
					require.Equal(
						t,
						property.DimCandidateStatusInAnalysis,
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
}
