package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
	"api5back/src/processing"

	"entgo.io/ent/dialect/sql"
)

func GetVacancySuggestions(
	ctx context.Context,
	client *ent.Client,
	pageRequest model.SuggestionsFilter,
) (*model.Page[model.Suggestion], error) {
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
		WithDimVacancy()

	if pageRequest.IDs != nil && len(*pageRequest.IDs) > 0 {
		query = query.
			Where(facthiringprocess.
				DimProcessIdIn(*pageRequest.IDs...))
	}

	page, pageSize, err := processing.ParsePageAndPageSize(
		pageRequest.Page,
		pageRequest.PageSize,
	)
	if err != nil {
		return nil, err
	}

	totalRecords, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	offset, numMaxPages := processing.ParseOffsetAndTotalPages(
		page,
		pageSize,
		totalRecords,
	)

	factHiringProcesses, err := query.
		Offset(offset).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var suggestions []model.Suggestion
	for _, fact := range factHiringProcesses {
		if fact.Edges.DimVacancy != nil {
			vacancy := fact.Edges.DimVacancy
			suggestions = append(suggestions, model.Suggestion{
				Id:    vacancy.DbId,
				Title: vacancy.Title,
			})
		}
	}

	return &model.Page[model.Suggestion]{
		Items:       suggestions,
		NumMaxPages: numMaxPages,
	}, nil
}
