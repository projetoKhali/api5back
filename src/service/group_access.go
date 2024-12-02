package service

import (
	"context"
	"errors"
	"fmt"

	"api5back/ent"
	"api5back/ent/department"
	"api5back/src/model"
)

func GetAccessGroupWithDepartments(
	ctx context.Context,
	client *ent.Client,
) ([]model.AccessGroup, error) {
	groups, err := client.
		AccessGroup.
		Query().
		WithDepartment().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var response []model.AccessGroup
	for _, group := range groups {
		var departments []model.Suggestion
		accessGroupDepartments, err := group.
			Edges.
			DepartmentOrErr()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve `Department` of `AccessGroup`: %w", err)
		}

		for _, dept := range accessGroupDepartments {
			departments = append(departments, model.Suggestion{
				Id:    dept.ID,
				Title: dept.Name,
			})
		}

		response = append(response, model.AccessGroup{
			Id:          group.ID,
			Name:        group.Name,
			Departments: departments,
		})
	}

	return response, nil
}

func CreateAccessGroup(
	ctx context.Context,
	client *ent.Client,
	request model.CreateAccessGroupRequest,
) (*model.AccessGroupCreated, error) {
	if request.Name == "" {
		return nil, errors.New("access group name cannot be empty")
	}

	departments, err := client.
		Department.
		Query().
		Where(department.IDIn(request.DepartmentIDs...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query departments: %w", err)
	}

	if len(departments) != len(request.DepartmentIDs) {
		return nil, errors.New("one or more department IDs do not exist")
	}

	group, err := client.
		AccessGroup.
		Create().
		SetName(request.Name).
		AddDepartment(departments...).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create access group: %w", err)
	}

	response := &model.AccessGroupCreated{
		ID:   group.ID,
		Name: group.Name,
	}

	return response, nil
}
