package service

import (
	"api5back/ent"
	"api5back/src/model"
	"context"
)

func ListDepartments(
	ctx context.Context,
	client *ent.Client,
) ([]model.Suggestion, error) {
	query := client.Department.Query()

	departments, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var response []model.Suggestion
	for _, dept := range departments {
		response = append(response, model.Suggestion{
			Id:    dept.ID,
			Title: dept.Name,
		})
	}

	return response, nil
}
