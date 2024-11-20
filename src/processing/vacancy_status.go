package processing

import (
	"fmt"

	"api5back/ent"
	"api5back/src/property"
)

type VacancyStatusSummary struct {
	Open      int `json:"open" default:"0"`
	Analyzing int `json:"analyzing" default:"0"`
	Closed    int `json:"closed" default:"0"`
}

func GenerateVacancyStatusSummary(
	data []*ent.FactHiringProcess,
) (VacancyStatusSummary, error) {
	countByStatus := make(map[property.DimVacancyStatus]int)

	for _, process := range data {
		vacancy, err := process.
			Edges.
			DimVacancyOrErr()
		if err != nil {
			return VacancyStatusSummary{}, fmt.Errorf(
				"`DimVacancy` of `FactHiringProcess` not found: %w",
				err,
			)
		}

		countByStatus[property.DimVacancyStatus(vacancy.Status+1)]++
	}

	return VacancyStatusSummary{
		Open:      countByStatus[property.DimVacancyStatusOpen],
		Analyzing: countByStatus[property.DimVacancyStatusInAnalysis],
		Closed:    countByStatus[property.DimVacancyStatusClosed],
	}, nil
}
