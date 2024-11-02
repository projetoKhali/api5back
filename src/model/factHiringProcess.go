package model

import "api5back/src/processing"

type MetricsData struct {
	VacancySummary    processing.VacancyStatusSummary      `json:"vacancyStatus"`
	CardInfos         processing.CardInfos                 `json:"cards"`
	AverageHiringTime processing.AverageHiringTimePerMonth `json:"averageHiringTime"`
}

type TableData struct {
	ProcessTitle      string   `json:"processTitle"`
	VacancyTitle      string   `json:"vacancyTitle"`
	NumPositions      int      `json:"numPositions"`
	NumCandidates     int      `json:"numCandidates"`
	CompetitionRate   *float32 `json:"competitionRate"`
	NumInterviewed    int      `json:"numInterviewed"`
	NumHired          int      `json:"numHired"`
	AverageHiringTime *float32 `json:"averageHiringTime"`
	NumFeedback       int      `json:"numFeedback"`
}
