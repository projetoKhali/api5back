package model

import "api5back/src/processing"

type DashboardMetrics struct {
	VacancySummary    processing.VacancyStatusSummary      `json:"vacancyStatus"`
	CardInfos         processing.CardInfos                 `json:"cards"`
	AverageHiringTime processing.AverageHiringTimePerMonth `json:"averageHiringTime"`
}

type DashboardTableRow struct {
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
