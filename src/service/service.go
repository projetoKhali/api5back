package service

import (
	"context"
	"fmt"

	"api5back/ent"
	"api5back/src/processing"
)

type MetricsService struct {
	dbClient *ent.Client
}
type MetricsData struct {
	VacancySummary    processing.VacancyStatusSummary      `json:"summary"`
	CardInfos         processing.CardInfos                 `json:"cardInfos"`
	AverageHiringTime processing.AverageHiringTimePerMonth `json:"avgHiringTime"`
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

	cardInfo, err := processing.ComputingCardInfo(hiringProcess)
	if err != nil {
		return metricsData, fmt.Errorf(
			"could not calculate `CardInfo` data: %w",
			err,
		)
	}
	metricsData.CardInfos = cardInfo

	vacancyInfo, err := processing.GenerateVacancyStatusSummary(hiringProcess)
	if err != nil {
		return metricsData, fmt.Errorf(
			"could not generate `VacancyStatus` summary: %w",
			err,
		)
	}
	metricsData.VacancySummary = vacancyInfo

	averageHiringTime, err := processing.GenerateAverageHiringTime(hiringProcess)
	if err != nil {
		return metricsData, fmt.Errorf(
			"could not generate `AvgHiringTime` data: %w",
			err,
		)
	}
	metricsData.AverageHiringTime = averageHiringTime

	return metricsData, nil
}
