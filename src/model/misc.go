package model

type DateRange struct {
	StartDate string `json:"startDate" form:"startDate" time_format:"2024-10-01T00:00:00Z"`
	EndDate   string `json:"endDate" form:"endDate" time_format:"2024-10-01T00:00:00Z"`
}

type Page[T any] struct {
	Items       []T `json:"items"`
	NumMaxPages int `json:"numMaxPages"`
}

type PageRequest struct {
	Page     *int `json:"page" default:"1"`
	PageSize *int `json:"pageSize" default:"10"`
}

type FactHiringProcessFilter struct {
	Recruiters    []int      `json:"recruiters"`
	Processes     []int      `json:"processes"`
	Vacancies     []int      `json:"vacancies"`
	DateRange     *DateRange `json:"dateRange"`
	ProcessStatus []int      `json:"processStatus"`
	VacancyStatus []int      `json:"vacancyStatus"`
	PageRequest
}

type SuggestionsFilter struct {
	IDs *[]int `json:"ids"`
	PageRequest
}
