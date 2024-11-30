package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimdepartment"
	"api5back/ent/dimprocess"
	"api5back/ent/dimuser"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
)

func GetUsers(
	ctx context.Context,
	client *ent.Client,
	departmentIDs *[]int,
) ([]model.Suggestion, error) {
	query := client.
		DimUser.
		Query()
	if departmentIDs != nil && len(*departmentIDs) > 0 {
		query =
			query.
				Where(
					dimuser.HasFactHiringProcessWith(
						facthiringprocess.HasDimProcessWith(
							dimprocess.HasDimDepartmentWith(
								dimdepartment.IDIn(*departmentIDs...),
							),
						),
					),
				)
	}
	users, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	uniqueUsers := make(map[int]model.Suggestion)

	for _, user := range users {
		uniqueUsers[user.ID] = model.Suggestion{
			Id:    user.DbId,
			Title: user.Name,
		}
	}

	result := make([]model.Suggestion, 0, len(uniqueUsers))
	for _, user := range uniqueUsers {
		result = append(result, user)
	}

	return result, nil
}
