package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimvacancy"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
	"api5back/src/pagination"
	"api5back/src/processing"

	"entgo.io/ent/dialect/sql"
)

func GetVacancySuggestions(
	ctx context.Context,
	client *ent.Client,
	pageRequest *model.SuggestionsFilter,
) (*model.Page[model.Suggestion], error) {
	query := client.
		FactHiringProcess.
		Query().
		WithDimVacancy().
		Order(
			facthiringprocess.ByDimVacancyField(
				dimvacancy.FieldDbId,
			),
			ent.Desc(facthiringprocess.FieldID),
		).
		Modify(func(s *sql.Selector) {
			s.Select("DISTINCT ON (t1.db_id) *")
		}).
		Clone()

	if pageRequest != nil && pageRequest.IDs != nil && len(*pageRequest.IDs) > 0 {
		query = query.
			Where(facthiringprocess.
				DimProcessIdIn(*pageRequest.IDs...))
	}

	page, pageSize, err := pagination.ParsePageRequest(pageRequest)
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
