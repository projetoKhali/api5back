package processing

import (
	"fmt"

	"api5back/ent"
	"api5back/src/property"
)

type AverageHiringTimePerMonth struct {
	January   float32 `json:"january"`
	February  float32 `json:"february"`
	March     float32 `json:"march"`
	April     float32 `json:"april"`
	May       float32 `json:"may"`
	June      float32 `json:"june"`
	July      float32 `json:"july"`
	August    float32 `json:"august"`
	September float32 `json:"september"`
	October   float32 `json:"october"`
	November  float32 `json:"november"`
	December  float32 `json:"december"`
}

type Month struct {
	TotalDurationInDays float64
	HiredCandidates     float64
}

func GenerateAverageHiringTime(
	data []*ent.FactHiringProcess,
) (AverageHiringTimePerMonth, error) {
	monthsValues := [12]Month{}

	for _, process := range data {
		candidates, err := process.Edges.HiringProcessCandidatesOrErr()
		if err != nil {
			return AverageHiringTimePerMonth{}, fmt.Errorf(
				"`HiringProcessCandidates` of `FactHiringProcess` not found: %w",
				err,
			)
		}

		for _, candidate := range candidates {
			if candidate.Status == property.HiringProcessCandidateStatusHired {
				interval := candidate.UpdatedAt.Time.Sub(candidate.ApplyDate.Time)
				intervalDays := interval.Hours() / 24
				monthIndex := candidate.UpdatedAt.Time.Month() - 1
				monthsValues[monthIndex].TotalDurationInDays += intervalDays
				monthsValues[monthIndex].HiredCandidates++
			}
		}
	}

	return AverageHiringTimePerMonth{
		January:   float32(monthsValues[0].TotalDurationInDays / monthsValues[0].HiredCandidates),
		February:  float32(monthsValues[1].TotalDurationInDays / monthsValues[1].HiredCandidates),
		March:     float32(monthsValues[2].TotalDurationInDays / monthsValues[2].HiredCandidates),
		April:     float32(monthsValues[3].TotalDurationInDays / monthsValues[3].HiredCandidates),
		May:       float32(monthsValues[4].TotalDurationInDays / monthsValues[4].HiredCandidates),
		June:      float32(monthsValues[5].TotalDurationInDays / monthsValues[5].HiredCandidates),
		July:      float32(monthsValues[6].TotalDurationInDays / monthsValues[6].HiredCandidates),
		August:    float32(monthsValues[7].TotalDurationInDays / monthsValues[7].HiredCandidates),
		September: float32(monthsValues[8].TotalDurationInDays / monthsValues[8].HiredCandidates),
		October:   float32(monthsValues[9].TotalDurationInDays / monthsValues[9].HiredCandidates),
		November:  float32(monthsValues[10].TotalDurationInDays / monthsValues[10].HiredCandidates),
		December:  float32(monthsValues[11].TotalDurationInDays / monthsValues[11].HiredCandidates),
	}, nil
}
