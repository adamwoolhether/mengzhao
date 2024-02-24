package db

import (
	"context"

	"github.com/google/uuid"

	"mengzhao/types"
)

func GetAccountByID(ctx context.Context, id uuid.UUID) (types.Account, error) {
	var account types.Account

	err := Bun.NewSelect().Model(&account).Where("user_id =?", id).Scan(ctx)
	if err != nil {
		return types.Account{}, err
	}

	return account, nil
}

func CreateAccount(ctx context.Context, account *types.Account) error {
	_, err := Bun.NewInsert().Model(account).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAccount(ctx context.Context, account *types.Account) error {
	_, err := Bun.NewUpdate().Model(account).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
