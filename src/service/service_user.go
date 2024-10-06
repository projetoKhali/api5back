package service

import (
	"context"

	"api5back/ent"
	"api5back/ent/dimuser"
)

func GetUsers(
	ctx context.Context,
	client *ent.Client,
) ([]*ent.DimUser, error) {
	users, err := client.DimUser.
		Query().
		Select(dimuser.FieldID, dimuser.FieldName).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
