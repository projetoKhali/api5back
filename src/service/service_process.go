package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimprocess"
)

func ListHiringProcesses(
	ctx context.Context,
	client *ent.Client,
	userIDs []int,
) ([]*ent.DimProcess, error) {
	query := client.DimProcess.Query()

	if len(userIDs) > 0 {
		query = query.Where(dimprocess.DimUsrIdIn(userIDs...))
	}

	processes, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return processes, nil
}
