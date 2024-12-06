package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimdepartment"
	"api5back/ent/dimprocess"
	"api5back/ent/dimuser"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
	"api5back/src/pagination"
	"api5back/src/processing"
)

func GetUserSuggestions(
	ctx context.Context,
	client *ent.Client,
	pageRequest *model.SuggestionsPageRequest,
) (*model.Page[model.Suggestion], error) {
	query := client.
		DimUser.
		Query()

	if pageRequest != nil && pageRequest.DepartmentIds != nil && len(*pageRequest.DepartmentIds) > 0 {
		query = query.
			Where(
				dimuser.HasFactHiringProcessWith(
					facthiringprocess.HasDimProcessWith(
						dimprocess.HasDimDepartmentWith(
							dimdepartment.IDIn(*pageRequest.DepartmentIds...),
						),
					),
				),
			)
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

	users, err := query.
		Offset(offset).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var suggestions []model.Suggestion
	for _, user := range users {
		suggestions = append(suggestions, model.Suggestion{
			Id:    user.DbId,
			Title: user.Name,
		})
	}

	return &model.Page[model.Suggestion]{
		Items:       suggestions,
		NumMaxPages: numMaxPages,
	}, nil
}
