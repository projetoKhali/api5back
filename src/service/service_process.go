package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimdepartment"
	"api5back/ent/dimprocess"
	"api5back/src/model"
)

func ListHiringProcesses(
	ctx context.Context,
	client *ent.Client,
	body model.BodySuggestion,
) ([]model.Suggestion, error) {
	query := client.
		DimProcess.
		Query()

	if body.DepartmentIds != nil && len(*body.DepartmentIds) > 0 {
		query =
			query.
				Where(
					dimprocess.HasDimDepartmentWith(
						dimdepartment.IDIn(*body.DepartmentIds...),
					),
				)
	}

	if body.FilterIds != nil && len(*body.FilterIds) > 0 {
		query = query.Where(dimprocess.DimUsrIdIn(*body.FilterIds...))
	}

	processes, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var response []model.Suggestion
	for _, process := range processes {
		response = append(response, model.Suggestion{
			Id:    process.DbId,
			Title: process.Title,
		})
	}

	return response, nil
}
