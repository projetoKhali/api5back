package processing

import (
	"fmt"
	"time"

	"api5back/ent"
	"api5back/src/property"
)

type CardInfos struct {
	Open                int `json:"open" default:"0"`
	InProgress          int `json:"inProgress" default:"0"`
	Closed              int `json:"closed" default:"0"`
	ApproachingDeadline int `json:"approachingDeadline" default:"0"`
	AverageHiringTime   int `json:"averageHiringTime" default:"0"`
}

func ComputingCardsInfo(
	factHiringProcesses []*ent.FactHiringProcess,
) (CardInfos, error) {
	if len(factHiringProcesses) == 0 {
		return CardInfos{}, nil
	}

	countByStatus := make(map[property.DimProcessStatus]int)
	approachingDeadline := 0
	totalHiringTime := 0.0
	totalCandidates := 0

	for _, factHiringProcess := range factHiringProcesses {
		process, err := factHiringProcess.
			Edges.
			DimProcessOrErr()
		if err != nil {
			return CardInfos{}, fmt.Errorf(
				"error getting `dim_process` with ID %d of factHiringProcess with ID %d: %+v",
				factHiringProcess.DimProcessId,
				factHiringProcess.ID,
				err,
			)
		}

		countByStatus[property.DimProcessStatus(process.Status+1)]++

		totalDuration := process.
			FinishDate.
			Time.
			Sub(process.InitialDate.Time)

		twentyPercentDuration := totalDuration * 20 / 100

		if time.Until(process.FinishDate.Time) < twentyPercentDuration && process.Status == 1 {
			approachingDeadline++
		}

		vacancy, err := factHiringProcess.
			Edges.
			DimVacancyOrErr()
		if err != nil {
			return CardInfos{}, fmt.Errorf(
				"error getting `dim_vacancy` with ID %d of factHiringProcess with ID %d: %+v",
				factHiringProcess.DimVacancyId,
				factHiringProcess.ID,
				err,
			)
		}

		candidates, err := vacancy.
			Edges.
			HiringProcessCandidatesOrErr()
		if err != nil {
			return CardInfos{}, fmt.Errorf(
				"error getting `hiring_process_candidates` of `dim_vacancy` with ID %d of factHiringProcess with ID %d: %+v",
				factHiringProcess.DimVacancyId,
				factHiringProcess.ID,
				err,
			)
		}

		for _, candidate := range candidates {
			if candidate.Status == property.HiringProcessCandidateStatusHired {
				interval := candidate.UpdatedAt.Time.Sub(candidate.ApplyDate.Time)
				intervalDays := interval.Hours() / 24
				totalHiringTime += intervalDays
				totalCandidates++
			}
		}
	}

	return CardInfos{
		Open:                countByStatus[property.DimProcessStatusOpen],
		InProgress:          countByStatus[property.DimProcessStatusInProgress],
		Closed:              countByStatus[property.DimProcessStatusClosed],
		ApproachingDeadline: approachingDeadline,
		AverageHiringTime:   int(totalHiringTime / float64(totalCandidates)),
	}, nil
}
