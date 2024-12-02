package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimdepartment"
	"api5back/ent/dimprocess"
	"api5back/src/model"
	"api5back/src/pagination"
	"api5back/src/processing"

	"entgo.io/ent/dialect/sql"
)

func GetProcessSuggestions(
	ctx context.Context,
	client *ent.Client,
	pageRequest *model.SuggestionsFilter,
) (*model.Page[model.Suggestion], error) {
	query := client.
		DimProcess.
		Query().
		Order(
			ent.Desc(dimprocess.FieldDbId),
			ent.Desc(dimprocess.FieldID),
		).
		Modify(func(s *sql.Selector) {
			s.Select("DISTINCT ON (db_id) *")
		}).
		Clone()

	if pageRequest != nil {
		if pageRequest.DepartmentIds != nil && len(*pageRequest.DepartmentIds) > 0 {
			query = query.
				Where(
					dimprocess.HasDimDepartmentWith(
						dimdepartment.IDIn(*pageRequest.DepartmentIds...),
					),
				)
		}

		if pageRequest.IDs != nil && len(*pageRequest.IDs) > 0 {
			query = query.
				Where(dimprocess.
					DimUsrIdIn(*pageRequest.IDs...))
		}
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

	processes, err := query.
		Offset(offset).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var suggestions []model.Suggestion
	for _, process := range processes {
		suggestions = append(suggestions, model.Suggestion{
			Id:    process.DbId,
			Title: process.Title,
		})
	}

	return &model.Page[model.Suggestion]{
		Items:       suggestions,
		NumMaxPages: numMaxPages,
	}, nil
}
