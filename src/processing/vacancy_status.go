package processing

import (
	"fmt"

	"api5back/ent"
)

type VacancyStatusSummary struct {
	Open    int `json:"open"`
	Hired   int `json:"hired"`
	Expired int `json:"expired"`
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
			summary.Hired++
		case 2:
			summary.Expired++
		default:
			return summary, fmt.Errorf(
				"invalid vacancy status: %d",
				vacancy.Status,
			)
		}
	}

	return summary, nil
}
