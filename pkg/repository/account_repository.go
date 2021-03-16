package repository

import (
	"account/pkg/accountinfo/dto"
	"account/pkg/accountinfo/model"
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AccountRepository interface {
	CreateOrUpdateAccountInfo(ctx context.Context, accountInfo *model.AccountInfo) error
	GetAccountData(ctx context.Context, accountQuery *dto.AccountQuery) ([]model.AccountInfo, error)
}

type gormAccountRepository struct {
	db *gorm.DB
}

func (gar *gormAccountRepository) CreateOrUpdateAccountInfo(ctx context.Context, accountInfo *model.AccountInfo) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	db := gar.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"balance",
				"updated_at",
			}),
		}).Create(&accountInfo)
	if db.Error != nil {
		return fmt.Errorf("create account entry failed. error %w", db.Error)
	}

	return nil
}

func (gar *gormAccountRepository) GetAccountData(ctx context.Context, accountQuery *dto.AccountQuery) ([]model.AccountInfo, error) {
	var res []model.AccountInfo
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	query := gar.db.WithContext(ctx)
	if accountQuery != nil && accountQuery.UserID.String() != "" {
		query = query.Where("user_id = ?", accountQuery.UserID.String())
	}
	if accountQuery != nil && accountQuery.AccountID.String() != "" {
		query = query.Where("id = ?", accountQuery.AccountID.String())
	}
	db := query.Find(&res)
	if db.Error != nil {
		return nil, fmt.Errorf("get account data failed: %w", db.Error)
	}

	return res, nil
}

func NewAccountInfoRepository(db *gorm.DB) AccountRepository {
	return &gormAccountRepository{
		db: db,
	}
}
