package seeds

import (
	"context"
	"fmt"

	"api5back/ent"
)

// Função para popular os dados no banco
func DataRelational(client *ent.Client) error {
	ctx := context.Background()

	departments := []struct {
		ID          int
		Name        string
		Description string
	}{
		{Name: "Marketing", Description: "MKT"},
		{Name: "RH", Description: "RH"},
		{Name: "CX", Description: "CX"},
		{Name: "ADM", Description: "Admin"},
		{Name: "Vendas", Description: "Vendas"},
	}

	departmentIDs := []int{}
	for _, d := range departments {
		dept, err := client.
			Department.
			Create().
			SetName(d.Name).
			SetDescription(d.Description).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create department %s: %v", d.Name, err)
		}
		departmentIDs = append(departmentIDs, dept.ID)
	}

	// Criar os grupos de acesso
	groups := []struct {
		Name     string
		GroupIDs []int
	}{
		{Name: "ADM", GroupIDs: []int{1, 3}},
		{Name: "RH", GroupIDs: []int{2, 4}},
		{Name: "Comercial", GroupIDs: []int{1, 5}},
		{Name: "Gestão", GroupIDs: []int{1, 2, 4}},
		{Name: "Vendas", GroupIDs: []int{3, 5}},
	}

	accessGroupIds := []int{}
	for _, g := range groups {
		group, err := client.
			AccessGroup.
			Create().
			SetName(g.Name).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create Group Access %s: %v", g.Name, err)
		}
		accessGroupIds = append(accessGroupIds, group.ID)
	}

	relations := []struct {
		DepartmentID int
		GroupIDs     []int
	}{
		{DepartmentID: 1, GroupIDs: []int{1, 3}},
		{DepartmentID: 2, GroupIDs: []int{2, 4}},
		{DepartmentID: 3, GroupIDs: []int{1, 5}},
		{DepartmentID: 4, GroupIDs: []int{1, 2, 4}},
		{DepartmentID: 5, GroupIDs: []int{3, 5}},
	}

	for _, rel := range relations {
		dept := departmentIDs[rel.DepartmentID-1]
		groupsIDs := []int{}
		for _, relationGroupID := range rel.GroupIDs {
			groupsIDs = append(groupsIDs, accessGroupIds[relationGroupID-1])
		}

		err := client.
			Department.
			UpdateOneID(dept).
			AddAccessGroupIDs(groupsIDs...).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create relation for department %d: %v", dept, err)
		}
	}

	users := []ent.Authentication{
		{Name: "Alice Santos", Email: "AliceSantos@gmail.com", Password: "password123", GroupId: 1},
		{Name: "Bob Ferreira", Email: "BobFerreira@gmail.com", Password: "password123", GroupId: 2},
		{Name: "Carla Mendes", Email: "CarlaMendes@gmail.com", Password: "password123", GroupId: 3},
		{Name: "David Costa", Email: "DavidCosta@gmail.com", Password: "password123", GroupId: 4},
		{Name: "Eva Lima", Email: "EvaLima@gmail.com", Password: "password123", GroupId: 5},
	}

	for _, user := range users {
		_, err := client.
			Authentication.
			Create().
			SetName(user.Name).
			SetEmail(user.Email).
			SetPassword(user.Password).
			SetAccessGroupID(user.GroupId).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", user.Name, err)
		}
	}

	return nil
}
