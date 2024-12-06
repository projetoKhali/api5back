package service

import (
	"context"
	"fmt"

	"api5back/ent"
	"api5back/ent/dimdepartment"
	"api5back/ent/dimprocess"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
	"api5back/src/pagination"
	"api5back/src/processing"
)

func GetVacancySuggestions(
	ctx context.Context,
	client *ent.Client,
	pageRequest *model.SuggestionsFilter,
) (*model.Page[model.Suggestion], error) {
	query := client.
		FactHiringProcess.
		Query().
		WithDimVacancy()

	if pageRequest != nil {
		if pageRequest.DepartmentIds != nil && len(*pageRequest.DepartmentIds) > 0 {
			query = query.
				Where(
					facthiringprocess.HasDimProcessWith(
						dimprocess.HasDimDepartmentWith(
							dimdepartment.IDIn(*pageRequest.DepartmentIds...),
						),
					),
				)
		}

		if pageRequest.IDs != nil && len(*pageRequest.IDs) > 0 {
			query = query.
				Where(facthiringprocess.
					DimProcessIdIn(*pageRequest.IDs...))
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

	fmt.Printf("pageSize: %+v | totalRecords: %+v\n", pageSize, totalRecords)
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
