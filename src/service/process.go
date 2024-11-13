package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimprocess"
	"api5back/src/model"
)

func ListHiringProcesses(
	ctx context.Context,
	client *ent.Client,
	userIDs *[]int,
) ([]model.Suggestion, error) {
	query := client.
		DimProcess.
		Query()

	// Verificar se o array de userIDs não é nil e se tem elementos
	if userIDs != nil && len(*userIDs) > 0 {
		query = query.Where(dimprocess.DimUsrIdIn(*userIDs...))
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
