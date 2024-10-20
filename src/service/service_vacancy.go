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

type UniqueVacancy struct {
	ID    int
	Title string
}

func (vs *VacancyService) GetVacancySuggestions(ctx context.Context, processesIds []int) ([]UniqueVacancy, error) {
	query := vs.dbClient.FactHiringProcess.Query()

	if len(processesIds) > 0 {
		query = query.Where(facthiringprocess.DimProcessIdIn(processesIds...))
		query = query.WithDimVacancy()

	}

	vacancies, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	uniqueVacancies := make(map[int]UniqueVacancy)

	for _, fact := range vacancies {
		if fact.Edges.DimVacancy != nil {
			vacancy := fact.Edges.DimVacancy
			uniqueVacancies[vacancy.ID] = UniqueVacancy{
				ID:    vacancy.ID,
				Title: vacancy.Title,
			}
		}
	}

	// Convert the map to a slice
	result := make([]UniqueVacancy, 0, len(uniqueVacancies))
	for _, v := range uniqueVacancies {
		result = append(result, v)
	}

	return result, nil
}
