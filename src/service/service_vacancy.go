package service

import (
	"api5back/ent"
	"api5back/ent/facthiringprocess"
	"context"
)

type VacancyService struct {
	dbClient *ent.Client
}

func NewVacancyService(client *ent.Client) *VacancyService {
	return &VacancyService{dbClient: client}
}

func (vs *VacancyService) GetVacancySuggestions(ctx context.Context, processesIds []int) ([]*ent.FactHiringProcess, error) {
	query := vs.dbClient.FactHiringProcess.Query()

	if len(processesIds) > 0 {
		query = query.Where(facthiringprocess.DimProcessIdIn(processesIds...))
		query = query.WithDimVacancy()

	}

	vacancies, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return vacancies, nil
}
