package service

import (
	"api5back/ent"
	"api5back/ent/department"
	"api5back/src/model"
	"context"
	"errors"
	"fmt"
)

type GroupAcessReturn struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`
	Departments []model.Suggestion `json:"departments"`
}

type CreateGroupAcessRequest struct {
	Name          string `json:"name" binding:"required"`
	DepartmentIDs []int  `json:"departments" binding:"required"`
}

type CreateGroupAcessResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetGroupAcessWithDepartments(
	ctx context.Context,
	client *ent.Client,
) ([]GroupAcessReturn, error) {
	groups, err := client.GroupAcess.Query().
		WithDepartment().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var response []GroupAcessReturn
	for _, group := range groups {
		var departments []model.Suggestion
		for _, dept := range group.Edges.Department {
			departments = append(departments, model.Suggestion{
				Id:    dept.ID,
				Title: dept.Name,
			})
		}

		response = append(response, GroupAcessReturn{
			Id:          group.ID,
			Name:        group.Name,
			Departments: departments,
		})
	}

	return response, nil
}

func CreateGroupAcess(
	ctx context.Context,
	client *ent.Client,
	request CreateGroupAcessRequest,
) (*CreateGroupAcessResponse, error) {
	if request.Name == "" {
		return nil, errors.New("group name cannot be empty")
	}

	departments, err := client.Department.Query().
		Where(department.IDIn(request.DepartmentIDs...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query departments: %w", err)
	}

	if len(departments) != len(request.DepartmentIDs) {
		return nil, errors.New("one or more department IDs do not exist")
	}

	group, err := client.GroupAcess.Create().
		SetName(request.Name).
		AddDepartment(departments...).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create group access: %w", err)
	}

	response := &CreateGroupAcessResponse{
		ID:   group.ID,
		Name: group.Name,
	}

	return response, nil
}
