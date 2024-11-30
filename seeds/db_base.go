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
		{ID: 1, Name: "Marketing", Description: "MKT"},
		{ID: 2, Name: "RH", Description: "RH"},
		{ID: 3, Name: "CX", Description: "CX"},
		{ID: 4, Name: "ADM", Description: "Admin"},
		{ID: 5, Name: "Vendas", Description: "Vendas"},
	}

	departmentMap := make(map[int]*ent.Department)
	for _, d := range departments {
		dept, err := client.Department.Create().
			SetID(d.ID).
			SetName(d.Name).
			SetDescription(d.Description).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create department %s: %v", d.Name, err)
		}
		departmentMap[d.ID] = dept
	}

	// Criar os grupos de acesso
	groupAccesses := []struct {
		ID   int
		Name string
	}{
		{ID: 1, Name: "ADM"},
		{ID: 2, Name: "RH"},
		{ID: 3, Name: "Comercial"},
		{ID: 4, Name: "Gestão"},
		{ID: 5, Name: "Vendas"},
	}

	groupAccessMap := make(map[int]*ent.GroupAcess)
	for _, g := range groupAccesses {
		group, err := client.GroupAcess.Create().
			SetID(g.ID).
			SetName(g.Name).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create Group Access %s: %v", g.Name, err)
		}
		groupAccessMap[g.ID] = group
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
		dept := departmentMap[rel.DepartmentID]
		var groups []*ent.GroupAcess
		for _, gid := range rel.GroupIDs {
			groups = append(groups, groupAccessMap[gid])
		}

		err := client.Department.UpdateOne(dept).
			AddGroupAcess(groups...).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create relation for department %s: %v", dept.Name, err)
		}
	}

	users := []ent.Authentication{
		{ID: 1, Name: "Alice Santos", Email: "AliceSantos@gmail.com", Password: "password123", GroupId: 1},
		{ID: 2, Name: "Bob Ferreira", Email: "BobFerreira@gmail.com", Password: "password123", GroupId: 2},
		{ID: 3, Name: "Carla Mendes", Email: "CarlaMendes@gmail.com", Password: "password123", GroupId: 3},
		{ID: 4, Name: "David Costa", Email: "DavidCosta@gmail.com", Password: "password123", GroupId: 4},
		{ID: 5, Name: "Eva Lima", Email: "EvaLima@gmail.com", Password: "password123", GroupId: 5},
	}

	for _, user := range users {
		_, err := client.Authentication.Create().
			SetID(user.ID).
			SetName(user.Name).
			SetEmail(user.Email).
			SetPassword(user.Password).
			SetGroupAcessID(user.GroupId).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", user.Name, err)
		}
	}

	return nil
}
