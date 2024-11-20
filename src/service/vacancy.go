package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"

	"entgo.io/ent/dialect/sql"
)

func GetVacancySuggestions(
	ctx context.Context,
	client *ent.Client,
	processesIds *[]int,
) ([]model.Suggestion, error) {
	query := client.
		FactHiringProcess.
		Query().
		Modify(func(s *sql.Selector) {
			s.Select("DISTINCT ON (t1.db_id) *")
			s.Join(sql.Table("dim_vacancy t1")).On(
				s.C("dim_vacancy_id"), sql.Table("t1").C("id"),
			)
		}).
		Order(ent.Desc(facthiringprocess.FieldID)).
		WithDimVacancy() // Com WithDimVacancy sempre incluído

	// Se processesIds não for nil e não estiver vazio, aplicamos o filtro
	if processesIds != nil && len(*processesIds) > 0 {
		query = query.Where(facthiringprocess.DimProcessIdIn(*processesIds...))
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
				Id:    vacancy.ID,
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
