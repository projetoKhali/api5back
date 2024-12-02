package model

type AccessGroup struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Departments []Suggestion `json:"departments"`
}

type CreateAccessGroupRequest struct {
	Name          string `json:"name" binding:"required"`
	DepartmentIDs []int  `json:"departments" binding:"required"`
}

type AccessGroupCreated struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
