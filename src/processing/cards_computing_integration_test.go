//go:build integration
// +build integration

package processing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"api5back/ent"
	"api5back/seeds"
	"api5back/src/database"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestData(
	client *ent.Client,
) (
	[]*ent.DimProcessCreate,
	[]*ent.FactHiringProcessCreate,
	[]*ent.DimCandidateCreate,
) {
	var processes []*ent.DimProcessCreate
	var factHiringProcesses []*ent.FactHiringProcessCreate
	var candidates []*ent.DimCandidateCreate

	for datePairIndex, datePair := range [3][2]string{
		{"2022-01-19", "2022-01-21"},
		{"2022-02-14", "2022-02-21"},
		{"2022-02-28", "2022-03-08"},
	} {
		initialDate, _ := time.Parse("2006-01-02", datePair[0])
		finishDate, _ := time.Parse("2006-01-02", datePair[1])

		processes = append(processes, client.
			DimProcess.
			Create().
			SetDbId(datePairIndex).
			SetTitle("Process "+datePair[0]).
			SetInitialDate(&pgtype.Date{Time: initialDate, Valid: true}).
			SetFinishDate(&pgtype.Date{Time: finishDate, Valid: true}).
			SetStatus(1).
			SetDimUsrId(1))

		factHiringProcesses = append(factHiringProcesses, client.
			FactHiringProcess.
			Create().
			SetDimProcessID(datePairIndex).
			SetDimVacancyID(datePairIndex).
			SetDimUserID(1).
			SetDimDatetimeID(1),
		)

		candidates = append(candidates, client.
			DimCandidate.
			Create().
			SetName(fmt.Sprintf("Candidate [%d]", datePairIndex)).
			SetEmail("can@didate.com").
			SetPhone("123456789").
			SetScore(100).
			SetDbId(datePairIndex).
			SetApplyDate(&pgtype.Date{Time: initialDate, Valid: true}).
			SetUpdatedAt(&pgtype.Date{Time: finishDate, Valid: true}).
			SetDimVacancyDbId(datePairIndex).
			SetStatus(1))

	}

	return processes, factHiringProcesses, candidates
}

func TestComputingCardInfo(t *testing.T) {
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

	if testResult := t.Run("Test ComputingCardsInfo", func(t *testing.T) {
		factHiringProcesses, err := intEnv.
			Client.
			FactHiringProcess.
			Query().
			WithDimProcess().
			WithDimVacancy(func(q *ent.DimVacancyQuery) {
				q.WithDimCandidates()
			}).
			All(ctx)

		require.NoError(t, err)
		require.NotEmpty(t, factHiringProcesses)

		cardInfos, err := ComputingCardsInfo(factHiringProcesses)
		require.NoError(t, err)

		assert.Equal(t, 9, cardInfos.Open)
		assert.Equal(t, 1, cardInfos.InProgress)
		assert.Equal(t, 0, cardInfos.Closed)
		assert.Equal(t, 1, cardInfos.ApproachingDeadline)
		assert.Equal(t, 14, cardInfos.AverageHiringTime)
	}); !testResult {
		t.Fatalf("Failed to query FactHiringProcess: %v", err)
	}
}

func TestComputingCardInfo_EmptyData(t *testing.T) {
	// Chama a função com uma lista vazia
	cardInfos, err := ComputingCardsInfo([]*ent.FactHiringProcess{})

	// Verifica se não houve erro
	assert.NoError(t, err)

	// Verifica se os valores do cardInfos são os valores padrão (zero)
	assert.Equal(t, CardInfos{}, cardInfos)
}
