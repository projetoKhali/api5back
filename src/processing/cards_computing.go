package processing

import (
	"fmt"
	"time"

	"api5back/ent"
)

// CardInfos is a struct that contains the card computing data
type CardInfos struct {
	OpenProcesses                int `json:"openProcess"`
	ExpiredProcesses             int `json:"expirededProcess"`
	ApproachingDeadlineProcesses int `json:"approachingDeadlineProcess"`
	CloseProcesses               int `json:"closeProcess"`
	AverageHiringTime            int `json:"averageHiringTime"`
}

func ComputingCardInfo(
	hiringData []*ent.FactHiringProcess,
) (CardInfos, error) {
	if len(hiringData) == 0 {
		return CardInfos{}, nil
	}

	var cardInfos CardInfos
	cardInfos.OpenProcesses = 0
	cardInfos.ExpiredProcesses = 0
	cardInfos.ApproachingDeadlineProcesses = 0
	cardInfos.CloseProcesses = 0
	cardInfos.AverageHiringTime = 0
	totalHiringTime := 0

	for _, hiring := range hiringData {
		process, err := hiring.Edges.DimProcessOrErr()
		if err != nil {
			return cardInfos, fmt.Errorf("error getting process data: %v", err)
		}

		switch process.Status {
		case 1:
			cardInfos.OpenProcesses++
		case 2:
			cardInfos.ExpiredProcesses++
		case 3:
			cardInfos.CloseProcesses++
		}
		totalDuration := process.FinishDate.Sub(process.InitialDate)
		twentyPercentDuration := totalDuration * 20 / 100
		if time.Until(process.FinishDate) < twentyPercentDuration && process.Status == 1 {
			cardInfos.ApproachingDeadlineProcesses++
		}
		totalHiringTime += hiring.MetSumDurationHiringProces
	}

	cardInfos.AverageHiringTime = totalHiringTime / len(hiringData)

	return cardInfos, nil
}
