package service

import (
	"context"
	"errors"
	"fmt"

	"api5back/ent"
	"api5back/ent/authentication"
	"api5back/src/model"
)

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Group string `json:"group"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Group       string             `json:"group"`
	Departments []model.Suggestion `json:"departments"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	GroupID  int    `json:"group_id" binding:"required"`
}

type CreateUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetAllUsers(
	ctx context.Context,
	client *ent.Client,
) ([]UserResponse, error) {
	users, err := client.Authentication.Query().
		WithGroupAcess().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			Name:  user.Name,
			Email: user.Email,
			Group: user.Edges.GroupAcess.Name,
		})
	}

	return response, nil
}

func Login(
	ctx context.Context,
	client *ent.Client,
	request LoginRequest,
) (*LoginResponse, error) {
	user, err := client.Authentication.Query().
		Where(
			authentication.Email(request.Email),
			authentication.Password(request.Password),
		).
		WithGroupAcess(func(gaq *ent.GroupAcessQuery) {
			gaq.WithDepartment()
		}).
		Only(ctx)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	var departments []model.Suggestion
	for _, dept := range user.Edges.GroupAcess.Edges.Department {
		departments = append(departments, model.Suggestion{
			Id:    dept.ID,
			Title: dept.Name,
		})
	}

	response := &LoginResponse{
		Name:        user.Name,
		Email:       user.Email,
		Group:       user.Edges.GroupAcess.Name,
		Departments: departments,
	}

	return response, nil
}

func CreateUser(
	ctx context.Context,
	client *ent.Client,
	request CreateUserRequest,
) (*CreateUserResponse, error) {
	if request.Name == "" || request.Email == "" || request.Password == "" {
		return nil, errors.New("name, email, and password cannot be empty")
	}

	group, err := client.GroupAcess.Get(ctx, request.GroupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %w", err)
	}

	user, err := client.Authentication.Create().
		SetName(request.Name).
		SetEmail(request.Email).
		SetPassword(request.Password).
		SetGroupId(group.ID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := &CreateUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return response, nil
}
