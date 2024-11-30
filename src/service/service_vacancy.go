package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimdepartment"
	"api5back/ent/dimprocess"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
)

func GetVacancySuggestions(
	ctx context.Context,
	client *ent.Client,
	body model.BodySuggestion,
) ([]model.Suggestion, error) {
	query := client.
		FactHiringProcess.
		Query().
		WithDimVacancy()

	if body.DepartmentIds != nil && len(*body.DepartmentIds) > 0 {
		query =
			query.
				Where(
					facthiringprocess.HasDimProcessWith( // Note que o escopo do predicado agora é "facthiringprocess"
						dimprocess.HasDimDepartmentWith(
							dimdepartment.IDIn(*body.DepartmentIds...),
						),
					),
				)
	}

	// Se processesIds não for nil e não estiver vazio, aplicamos o filtro
	if body.FilterIds != nil && len(*body.FilterIds) > 0 {
		query = query.
			Where(facthiringprocess.
				DimProcessIdIn(*body.FilterIds...))
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
