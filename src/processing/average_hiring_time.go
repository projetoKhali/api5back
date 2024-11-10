package processing

import (
	"api5back/ent"
	"api5back/src/property"
	"fmt"
	"reflect"
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

func GenerateAverageHiringTimePerMonth(
	data []*ent.FactHiringProcess,
) (AverageHiringTimePerMonth, error) {
	monthsValues := [12]Month{}

	for _, process := range data {
		candidates, err := process.
			Edges.
			HiringProcessCandidatesOrErr()
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

	result := AverageHiringTimePerMonth{}
	resultValue := reflect.ValueOf(&result).Elem()

	for i := 0; i < len(monthsValues); i++ {
		if monthsValues[i].HiredCandidates > 0 {
			fieldName := resultValue.Type().Field(i).Name

			duration := monthsValues[i].TotalDurationInDays
			candidates := monthsValues[i].HiredCandidates

			if duration == 0 || candidates == 0 {
				resultValue.FieldByName(fieldName).SetFloat(0)
				continue
			}

			avg := float32(duration / candidates)
			resultValue.FieldByName(fieldName).SetFloat(float64(avg))
		}
	}

	return result, nil
}

func GenerateAverageHiringTimePerFactHiringProcess(
	fact_hiring_process *ent.FactHiringProcess,
) (float32, error) {
	candidates, err := fact_hiring_process.Edges.HiringProcessCandidatesOrErr()
	if err != nil {
		return 0, fmt.Errorf(
			"`HiringProcessCandidates` of `FactHiringProcess` not found: %w",
			err,
		)
	}

	hiredCandidates := 0.0
	days := 0.0

	for _, candidate := range candidates {
		if candidate.Status == property.HiringProcessCandidateStatusHired {
			interval := candidate.UpdatedAt.Time.Sub(candidate.ApplyDate.Time)
			intervalDays := interval.Hours() / 24
			hiredCandidates += 1
			days += intervalDays

		}
	}

	if hiredCandidates == 0 {
		return 0, fmt.Errorf(
			"`No hired candidates found: %w",
			err,
		)
	}
	result := float32(days / hiredCandidates)

	return result, nil
}
