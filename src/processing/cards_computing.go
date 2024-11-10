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

func ComputingCardInfo(
	factHiringProcesses []*ent.FactHiringProcess,
) (CardInfos, error) {
	if len(factHiringProcesses) == 0 {
		return CardInfos{}, nil
	}

	countByStatus := make(map[property.DimProcessStatus]int)
	approachingDeadline := 0
	totalHiringTime := 0

	for _, factHiringProcess := range factHiringProcesses {
		process, err := factHiringProcess.
			Edges.
			DimProcessOrErr()
		if err != nil {
			return CardInfos{}, fmt.Errorf("error getting `dim_process` of factHiringProcess: %+v", err)
		}

		countByStatus[process.Status]++

		totalDuration := process.
			FinishDate.
			Time.
			Sub(process.InitialDate.Time)

		twentyPercentDuration := totalDuration * 20 / 100

		if time.Until(process.FinishDate.Time) < twentyPercentDuration && process.Status == 1 {
			approachingDeadline++
		}

		totalHiringTime += factHiringProcess.MetSumDurationHiringProces
	}

	return CardInfos{
		Open:                countByStatus[property.DimProcessStatusOpen],
		InProgress:          countByStatus[property.DimProcessStatusInProgress],
		Closed:              countByStatus[property.DimProcessStatusClosed],
		ApproachingDeadline: approachingDeadline,
		AverageHiringTime:   totalHiringTime / len(factHiringProcesses),
	}, nil
}
