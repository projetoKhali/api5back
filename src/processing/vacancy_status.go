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
	factHiringProcesses []*ent.FactHiringProcess,
) (VacancyStatusSummary, error) {
	countByStatus := make(map[property.DimVacancyStatus]int)

	for _, factHiringProcess := range factHiringProcesses {
		vacancy, err := factHiringProcess.
			Edges.
			DimVacancyOrErr()
		if err != nil {
			return VacancyStatusSummary{}, fmt.Errorf(
				"`DimVacancy` with ID %d of `FactHiringProcess` with ID %d not found: %w",
				factHiringProcess.DimVacancyId,
				factHiringProcess.ID,
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
