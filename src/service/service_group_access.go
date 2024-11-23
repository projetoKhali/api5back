package service

import (
	"api5back/ent"
	"api5back/ent/department"
	"api5back/src/model"
	"context"
	"errors"
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

// GetGroupAcessWithDepartments retorna grupos de acesso com seus departamentos
func GetGroupAcessWithDepartments(
	ctx context.Context,
	client *ent.Client,
) ([]GroupAcessReturn, error) {
	// Consulta todos os grupos de acesso e os departamentos relacionados
	groups, err := client.GroupAcess.Query().
		WithDepartment(). // Inclui os departamentos relacionados
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Mapeia o resultado para o formato esperado
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
) (*ent.GroupAcess, error) {
	// Verifica se o nome foi fornecido
	if request.Name == "" {
		return nil, errors.New("group name cannot be empty")
	}

	// Verifica se os departamentos existem
	departments, err := client.Department.Query().
		Where(department.IDIn(request.DepartmentIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Cria o grupo de acesso e associa os departamentos
	group, err := client.GroupAcess.Create().
		SetName(request.Name).
		AddDepartment(departments...).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return group, nil
}
