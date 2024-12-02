package model

// DateRange represents a date range with a start and end date.
type DateRange struct {
	StartDate string `json:"startDate" form:"startDate" time_format:"2024-10-01" default:""`
	EndDate   string `json:"endDate" form:"endDate" time_format:"2024-10-01" default:""`
}

// Page represents a page of items in paginated responses.
type Page[T any] struct {
	Items       []T `json:"items"`
	NumMaxPages int `json:"numMaxPages"`
}

// PageRequest is the base type of request for a page of items.
type PageRequest struct {
	Page     *int `json:"page" default:"1"`
	PageSize *int `json:"pageSize" default:"10"`
}

func (pr *PageRequest) GetPageRequest() *PageRequest {
	return pr
}

// allows for any type that implements GetPageRequest
// to be accepted by ParsePageRequest function
type PageRequester interface {
	GetPageRequest() *PageRequest
}

// FactHiringProcessFilter represents a filter for querying FactHiringProcess entities.
type FactHiringProcessFilter struct {
	Recruiters    []int      `json:"recruiters"`
	Processes     []int      `json:"processes"`
	Vacancies     []int      `json:"vacancies"`
	DateRange     *DateRange `json:"dateRange"`
	ProcessStatus []int      `json:"processStatus"`
	VacancyStatus []int      `json:"vacancyStatus"`
	AccessGroups  []int      `json:"accessGroup"`
	*PageRequest
}

// SuggestionsPageRequest represents a paginated query for suggestions.
// It includes a list of department IDs for filtering.
type SuggestionsPageRequest struct {
	DepartmentIds *[]int `json:"departments"`
	*PageRequest
}

func (s *SuggestionsPageRequest) GetPageRequest() *PageRequest {
	if s == nil {
		return nil
	}
	return s.PageRequest
}

// SuggestionsFilter represents a filter for querying suggestions.
// It includes a list of IDs for filtering, in addition to the
// department IDs and base page request.
type SuggestionsFilter struct {
	IDs *[]int `json:"ids"`
	SuggestionsPageRequest
}

func (s *SuggestionsFilter) GetPageRequest() *PageRequest {
	if s == nil {
		return nil
	}
	return s.PageRequest
}
