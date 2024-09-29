package processing

import (
	"api5back/ent"
	"api5back/src/property"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAverageHiringTime(t *testing.T) {

	layout := "2006-01-02 15:04:05"
	datesString := [24]string{
		"2022-01-19 00:00:00", "2022-01-21 00:00:00",
		"2022-02-14 00:00:00", "2022-02-21 00:00:00",
		"2022-02-28 00:00:00", "2022-03-08 00:00:00",
		"2022-03-24 00:00:00", "2022-03-31 00:00:00",
		"2022-03-22 00:00:00", "2022-03-30 00:00:00",
		"2022-03-27 00:00:00", "2022-03-28 00:00:00",
		"2022-03-14 00:00:00", "2022-03-21 00:00:00",
		"2022-04-20 00:00:00", "2022-04-24 00:00:00",
		"2022-08-19 00:00:00", "2022-08-22 00:00:00",
		"2022-09-15 00:00:00", "2022-09-23 00:00:00",
		"2022-10-17 00:00:00", "2022-10-18 00:00:00",
		"2022-12-08 00:00:00", "2022-12-09 00:00:00",
	}
	candidates := [12]ent.HiringProcessCandidate{}

	for i := 0; i < 12; i += 2 {
		applyDate := datesString[i]
		timeApplyDate, _ := time.Parse(layout, applyDate)
		pgtypeApplyDate := &pgtype.Date{}
		pgtypeApplyDate.Scan(timeApplyDate)

		updatedAt := datesString[i+1]
		timeUpdatedAt, _ := time.Parse(layout, updatedAt)
		pgtypeUpdatedAt := &pgtype.Date{}
		pgtypeUpdatedAt.Scan(timeUpdatedAt)
		candidates[i] = ent.HiringProcessCandidate{
			ApplyDate: pgtypeApplyDate,
			UpdatedAt: pgtypeUpdatedAt,
			Status:    property.HiringProcessCandidateStatus(i % 4),
		}
	}

	hiringData := []*ent.FactHiringProcess{
		{
			Edges: ent.FactHiringProcessEdges{
				HiringProcessCandidates: []*ent.HiringProcessCandidate{
					&candidates[0],
					&candidates[1],
					&candidates[2],
					&candidates[3],
				},
			},
		},
		{
			Edges: ent.FactHiringProcessEdges{
				HiringProcessCandidates: []*ent.HiringProcessCandidate{
					&candidates[4],
					&candidates[5],
					&candidates[6],
					&candidates[7],
				},
			},
		},
		{
			Edges: ent.FactHiringProcessEdges{
				HiringProcessCandidates: []*ent.HiringProcessCandidate{
					&candidates[8],
					&candidates[9],
					&candidates[10],
					&candidates[11],
				},
			},
		},
	}

	fmt.Printf("%+v", hiringData[0].Edges)

	months, err := GenerateAverageHiringTime(hiringData)
	assert.NoError(t, err)

	assert.Equal(t, float32(0), months.January)
	assert.Equal(t, float32(0), months.February)
	assert.Equal(t, float32(7.5), months.March)
	assert.Equal(t, float32(0), months.April)
	assert.Equal(t, float32(0), months.May)
	assert.Equal(t, float32(0), months.June)
	assert.Equal(t, float32(0), months.July)
	assert.Equal(t, float32(0), months.August)
	assert.Equal(t, float32(0), months.September)
	assert.Equal(t, float32(1), months.October)
	assert.Equal(t, float32(0), months.November)
	assert.Equal(t, float32(0), months.December)
}
