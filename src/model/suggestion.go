package model

type Suggestion struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type BodySuggestion struct {
	FilterIds     *[]int `json:"filterIds"`
	DepartmentIds *[]int `json:"departments"`
}
