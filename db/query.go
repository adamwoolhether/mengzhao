package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"mengzhao/types"
)

func GetImagesByUserID(ctx context.Context, userID uuid.UUID) ([]types.Image, error) {
	var images []types.Image

	err := Bun.NewSelect().
		Model(&images).
		Where("deleted = ?", false).
		Where("user_id =?", userID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func GetImageByID(ctx context.Context, imageID int) (types.Image, error) {
	var image types.Image

	err := Bun.NewSelect().Model(&image).Where("id =?", imageID).Scan(ctx)
	if err != nil {
		return image, err
	}

	return image, nil
}

func GetImagesByBatchID(ctx context.Context, batchID uuid.UUID) ([]types.Image, error) {
	var images []types.Image

	err := Bun.NewSelect().
		Model(&images).
		Where("batch_id =?", batchID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func CreateImage(ctx context.Context, tx bun.Tx, image *types.Image) error {
	_, err := tx.NewInsert().Model(image).Exec(ctx)
	if err != nil {
		return nil
	}

	return nil
}

func UpdateImage(ctx context.Context, tx bun.Tx, image *types.Image) error {
	_, err := tx.NewUpdate().Model(image).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

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
