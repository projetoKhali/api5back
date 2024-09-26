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
	processData []*ent.DimProcess,
	HiringData []*ent.FactHiringProcess,
) (CardInfos, error) {
	var cardInfos CardInfos

	err := getProcessCardInfo(processData, &cardInfos)
	if err != nil {
		return cardInfos, err
	}
	err = getAverageHiringTime(HiringData, &cardInfos)
	if err != nil {
		return cardInfos, err
	}

	return cardInfos, nil
}

func getProcessCardInfo(
	data []*ent.DimProcess,
	cardInfos *CardInfos,
) error {
	cardInfos.openProcess = 0
	cardInfos.expirededProcess = 0
	cardInfos.approachingDeadlineProcess = 0
	cardInfos.closeProcess = 0
	cardInfos.averageHiringTime = 0

	for _, process := range data {
		if process.Status == 1 {
			cardInfos.openProcess++
		}
		if process.Status == 2 {
			cardInfos.expirededProcess++
		}
		if process.Status == 3 {
			cardInfos.closeProcess++
		}

		// calculate the approaching deadline process
		totalDuration := process.FinishDate.Sub(process.InitialDate)
		twentyPercentDuration := totalDuration * 20 / 100
		if time.Until(process.FinishDate) < twentyPercentDuration && process.Status == 1 {
			cardInfos.approachingDeadlineProcess++
		}
	}
	return nil
}

func getAverageHiringTime(
	data []*ent.FactHiringProcess,
	cardInfos *CardInfos,
) error {
	if len(data) == 0 {
		return fmt.Errorf("empty hiring data")
	}

	totalHiringTime := 0
	for _, hiring := range data {
		totalHiringTime += hiring.MetSumDurationHiringProces
	}
	cardInfos.averageHiringTime = totalHiringTime / len(data)
	return nil
}
