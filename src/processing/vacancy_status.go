package processing

import (
	"fmt"

	"api5back/ent"
)

type VacancyStatusSummary struct {
	Open      int `json:"open"`
	Analyzing int `json:"analyzing"`
	Closed    int `json:"closed"`
}

func GenerateVacancyStatusSummary(
	data []*ent.FactHiringProcess,
) (VacancyStatusSummary, error) {
	var summary VacancyStatusSummary

	for _, process := range data {
		vacancy, err := process.Edges.DimVacancyOrErr()
		if err != nil {
			return summary, fmt.Errorf(
				"`DimVacancy` of `FactHiringProcess` not found: %w",
				err,
			)
		}

		switch vacancy.Status {
		case 0:
			summary.Open++
		case 1:
			summary.Analyzing++
		case 2:
			summary.Closed++
		default:
			return summary, fmt.Errorf(
				"invalid vacancy status: %d",
				vacancy.Status,
			)
		}
	}

	return summary, nil
}
