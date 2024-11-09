package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
)

func GetVacancySuggestions(
	ctx context.Context,
	client *ent.Client,
	processesIds *[]int,
) ([]model.Suggestion, error) {
	query := client.
		FactHiringProcess.
		Query().
		WithDimVacancy() // Com WithDimVacancy sempre incluído

	// Se processesIds não for nil e não estiver vazio, aplicamos o filtro
	if processesIds != nil && len(*processesIds) > 0 {
		query = query.
			Where(facthiringprocess.
				DimProcessIdIn(*processesIds...))
	}

	vacancies, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	uniqueVacancies := make(map[int]model.Suggestion)

	for _, fact := range vacancies {
		if fact.Edges.DimVacancy != nil {
			vacancy := fact.Edges.DimVacancy
			uniqueVacancies[vacancy.ID] = model.Suggestion{
				Id:    vacancy.DbId,
				Title: vacancy.Title,
			}
		}
	}

	// Convertendo o map para slice
	result := make([]model.Suggestion, 0, len(uniqueVacancies))
	for _, v := range uniqueVacancies {
		result = append(result, v)
	}

	return result, nil
}
