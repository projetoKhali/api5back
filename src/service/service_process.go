package service

import (
	"api5back/ent"
	"api5back/ent/dimprocess"
	"context"
)

type HiringProcessService struct {
	client *ent.Client
}

func NewHiringProcessService(client *ent.Client) *HiringProcessService {
	return &HiringProcessService{client: client}
}

func (s *HiringProcessService) ListHiringProcesses(ctx context.Context, userIDs []int) ([]*ent.DimProcess, error) {
	query := s.client.DimProcess.Query()

	if len(userIDs) > 0 {
		query = query.Where(dimprocess.DimUsrIdIn(userIDs...))
	}

	processes, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return processes, nil
}
