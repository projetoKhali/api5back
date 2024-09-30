package service

import (
	"api5back/ent"
	"api5back/src/processing"
	"context"
	"fmt"
)

type MetricsService struct {
	dbClient *ent.Client
}

type MetricsData struct {
	VacancySummary    processing.VacancyStatusSummary      `json:"vacancyStatus"`
	CardInfos         processing.CardInfos                 `json:"cards"`
	AverageHiringTime processing.AverageHiringTimePerMonth `json:"averageHiringTime"`
}

func NewMetricsService(dbclient *ent.Client) *MetricsService {
	return &MetricsService{dbClient: dbclient}
}

func (s *MetricsService) GetMetrics(ctx context.Context) (MetricsData, error) {
	var metricsData MetricsData

	hiringProcess, err := s.dbClient.
		FactHiringProcess.
		Query().
		WithDimVacancy().
		WithDimProcess().
		WithHiringProcessCandidates().
		All(ctx)
	if err != nil {
		return metricsData, fmt.Errorf(
			"could not retrieve `FactHiringProcess` data: %w",
			err,
		)
	}

	var errors []error

	cardInfo, err := processing.ComputingCardInfo(hiringProcess)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not calculate `CardInfo` data: %w",
			err,
		))
	}
	metricsData.CardInfos = cardInfo

	vacancyInfo, err := processing.GenerateVacancyStatusSummary(hiringProcess)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not generate `VacancyStatus` summary: %w",
			err,
		))
	}
	metricsData.VacancySummary = vacancyInfo

	averageHiringTime, err := processing.GenerateAverageHiringTime(hiringProcess)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not generate `AvgHiringTime` data: %w",
			err,
		))
	}
	metricsData.AverageHiringTime = averageHiringTime

	if len(errors) > 0 {
		return metricsData, fmt.Errorf(
			"failed to get metrics: %v",
			errors,
		)
	}

	return metricsData, nil
}
