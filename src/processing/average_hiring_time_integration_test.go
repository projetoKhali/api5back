//go:build integration
// +build integration

package processing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"api5back/ent"
	"api5back/ent/facthiringprocess"
	"api5back/seeds"
	"api5back/src/database"
	"api5back/src/property"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateTestData(
	client *ent.Client,
	ids []int,
) []*ent.HiringProcessCandidateCreate {
	datesString := [12][2]string{
		{"2022-01-19", "2022-01-21"},
		{"2022-02-14", "2022-02-21"},
		{"2022-02-28", "2022-03-08"},
		{"2022-03-24", "2022-03-31"},
		{"2022-03-22", "2022-03-30"},
		{"2022-03-27", "2022-03-28"},
		{"2022-03-14", "2022-03-21"},
		{"2022-04-20", "2022-04-24"},
		{"2022-08-19", "2022-08-22"},
		{"2022-09-15", "2022-09-23"},
		{"2022-10-17", "2022-10-18"},
		{"2022-12-08", "2022-12-09"},
	}
	candidates := []*ent.HiringProcessCandidateCreate{}

	for i, factID := range ids {
		factIndex := i * 8

		for j := 0; j < 4; j++ {
			dateIndex := factIndex + j

			stringApplyDate := datesString[dateIndex%12][0]
			timeApplyDate, _ := time.Parse(time.DateOnly, stringApplyDate)
			pgtypeApplyDate := &pgtype.Date{}
			pgtypeApplyDate.Scan(timeApplyDate)

			stringUpdatedAt := datesString[dateIndex%12][1]
			timeUpdatedAt, _ := time.Parse(time.DateOnly, stringUpdatedAt)
			pgtypeUpdatedAt := &pgtype.Date{}
			pgtypeUpdatedAt.Scan(timeUpdatedAt)

			candidates = append(candidates, client.
				HiringProcessCandidate.
				Create().
				SetDbId(j).
				SetName(fmt.Sprintf("Candidate[%d][%d]", i, j)).
				SetEmail("can@didate.com").
				SetPhone("123456789").
				SetScore(100).
				SetApplyDate(pgtypeApplyDate).
				SetUpdatedAt(pgtypeUpdatedAt).
				SetStatus(property.HiringProcessCandidateStatus(j)).
				SetFactHiringProcessID(factID),
			)
		}
	}

	return candidates
}

func TestAverageHiringTime(t *testing.T) {
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
	if testResult := t.Run("Query FactHiringProcess", func(t *testing.T) {
		factHiringProcesses, err = intEnv.
			Client.
			FactHiringProcess.
			Query().
			Limit(3).
			All(ctx)

		require.NoError(t, err)
		require.NotEmpty(t, factHiringProcesses)
	}); !testResult {
		t.Fatalf("Failed to query FactHiringProcess: %v", err)
	}

	var ids []int
	for _, factHiringProcess := range factHiringProcesses {
		ids = append(ids, factHiringProcess.ID)
	}

	candidates := CreateTestData(intEnv.Client, ids)
	if testResult := t.Run("Insert FactHiringProcessCandidates", func(t *testing.T) {
		for i := range ids {
			_, err = intEnv.
				Client.
				HiringProcessCandidate.
				CreateBulk(candidates[i*4 : (i*4)+4]...).
				Save(ctx)
			if err != nil {
				t.Fatalf("Failed to insert FactHiringProcessCandidates: %v", err)
			}
		}
	}); !testResult {
		t.Fatalf("Failed to insert FactHiringProcessCandidates: %v", err)
	}

	if testResult := t.Run("Select FactHiringProcess with HiringProcessCandidates", func(t *testing.T) {
		factHiringProcesses, err = intEnv.
			Client.
			FactHiringProcess.
			Query().
			WithHiringProcessCandidates().
			Where(facthiringprocess.IDIn(ids...)).
			All(ctx)

		require.NoError(t, err)
	}); !testResult {
		t.Fatalf("Failed to select FactHiringProcess with HiringProcessCandidates: %v", err)
	}

	if testResult := t.Run("Test GenerateAverageHiringTimePerMonth processing function", func(t *testing.T) {
		months, err := GenerateAverageHiringTimePerMonth(factHiringProcesses)
		require.NoError(t, err)

		assert.Equal(t, float32(0), months.January)
		assert.Equal(t, float32(0), months.February)
		assert.Equal(t, float32(8.5), months.March)
		assert.Equal(t, float32(0), months.April)
		assert.Equal(t, float32(0), months.May)
		assert.Equal(t, float32(14.5), months.June)
		assert.Equal(t, float32(13.5), months.July)
		assert.Equal(t, float32(0), months.August)
		assert.Equal(t, float32(13.666667), months.September)
		assert.Equal(t, float32(1), months.October)
		assert.Equal(t, float32(0), months.November)
		assert.Equal(t, float32(0), months.December)
	}); !testResult {
		t.Fatalf("Failed to test GenerateAverageHiringTimePerMonth processing function: %v", err)
	}
}
