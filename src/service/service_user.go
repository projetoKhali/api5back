package service

import (
	"api5back/ent"
	"api5back/ent/dimuser"
	"context"
)

type UserService struct {
	dwClient *ent.Client
}

func NewUserService(dwClient *ent.Client) *UserService {
	return &UserService{dwClient: dwClient}
}

func (s *UserService) GetUsers(ctx context.Context) ([]*ent.DimUser, error) {
	users, err := s.dwClient.DimUser.
		Query().
		Select(dimuser.FieldID, dimuser.FieldName).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
