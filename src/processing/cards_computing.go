package processing

import (
	"api5back/ent"
	"fmt"
	"time"
)

// CardInfos is a struct that contains the card computing data
type CardInfos struct {
	openProcess                int
	expirededProcess           int
	approachingDeadlineProcess int
	closeProcess               int
	averageHiringTime          int
}

func ComputingCardInfo(
	hiringData []*ent.FactHiringProcess,
) (CardInfos, error) {

	if len(hiringData) == 0 {
		return CardInfos{}, fmt.Errorf("the list is empty")
	}

	var cardInfos CardInfos
	cardInfos.openProcess = 0
	cardInfos.expirededProcess = 0
	cardInfos.approachingDeadlineProcess = 0
	cardInfos.closeProcess = 0
	cardInfos.averageHiringTime = 0
	totalHiringTime := 0

	for _, hiring := range hiringData {
		process, err := hiring.Edges.DimProcessOrErr()
		if err != nil {
			return cardInfos, fmt.Errorf("error getting process data: %v", err)
		}

		switch process.Status {
		case 1:
			cardInfos.openProcess++
		case 2:
			cardInfos.expirededProcess++
		case 3:
			cardInfos.closeProcess++
		}
		totalDuration := process.FinishDate.Sub(process.InitialDate)
		twentyPercentDuration := totalDuration * 20 / 100
		if time.Until(process.FinishDate) < twentyPercentDuration && process.Status == 1 {
			cardInfos.approachingDeadlineProcess++
		}
		totalHiringTime += hiring.MetSumDurationHiringProces
	}

	cardInfos.averageHiringTime = totalHiringTime / len(hiringData)

	return cardInfos, nil
}
