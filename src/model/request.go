package model

type DateRange struct {
	StartDate string `json:"startDate" form:"startDate" time_format:"2024-10-01T00:00:00Z"`
	EndDate   string `json:"endDate" form:"endDate" time_format:"2024-10-01T00:00:00Z"`
}

type FactHiringProcessFilter struct {
	Recruiters    []int      `json:"recruiters"`
	Processes     []int      `json:"processes"`
	Vacancies     []int      `json:"vacancies"`
	DateRange     *DateRange `json:"dateRange"`
	ProcessStatus []int      `json:"processStatus"`
	VacancyStatus []int      `json:"vacancyStatus"`
	Page          *int       `json:"page"`
	PageSize      *int       `json:"pageSize"`
}


