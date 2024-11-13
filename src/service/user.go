package service

import (
	"context"

	"api5back/ent"
	"api5back/src/model"
)

func GetUsers(
	ctx context.Context,
	client *ent.Client,
) ([]model.Suggestion, error) {
	users, err := client.DimUser.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var response []model.Suggestion
	for _, user := range users {
		response = append(response, model.Suggestion{
			Id:    user.DbId,
			Title: user.Name,
		})
	}

	return response, nil
}
