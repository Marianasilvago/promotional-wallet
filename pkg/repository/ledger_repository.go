package repository

import (
	"account/pkg/accountinfo/dto"
	"account/pkg/ledger/model"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type LedgerRepository interface {
	CreateLedgerEntry(ctx context.Context, ledgerEntry *model.Ledger) error
	GetEntries(ctx context.Context, query *dto.LogQuery) ([]*model.Ledger, error)
	//GetLedgerData(ctx context.Context, accountQuery *dto.LedgerQuery) ([]model.Ledger, error)
}

type gormLedgerRepository struct {
	db *gorm.DB
}

func (gbr *gormLedgerRepository) CreateLedgerEntry(ctx context.Context, ledgerEntry *model.Ledger) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	db := gbr.db.WithContext(ctx).Create(&ledgerEntry)
	if db.Error != nil {
		return fmt.Errorf("create ledger entry failed. error %w", db.Error)
	}

	return nil
}

func (gbr *gormLedgerRepository) GetEntries(ctx context.Context, logQuery *dto.LogQuery) ([]*model.Ledger, error) {
	entries := make([]*model.Ledger, 0)

	dbQuery := gbr.db.WithContext(ctx)
	if logQuery != nil {
		dbQuery.Where("activity = ?", logQuery.ActivityType)
	}
	db := dbQuery.Find(&entries)
	if db.Error != nil {
		return nil, fmt.Errorf("fetching ledger entry failed. error %w", db.Error)
	}

	return entries, nil
}

//
//func (gbr *gormLedgerRepository) GetLedgerData(ctx context.Context, accountQuery *dto.LedgerQuery) ([]model.Ledger, error) {
//	var res []model.Ledger
//	ctx, cancel := context.WithTimeout(ctx, time.Second)
//	defer cancel()
//
//	db := gbr.db.WithContext(ctx).Where("fsym in ? and tsym in ?", accountQuery.Fsyms, accountQuery.Tsyms).Find(&res)
//	if db.Error != nil {
//		return nil, fmt.Errorf("get account data failed: %w", db.Error)
//	}
//
//	return res, nil
//}

func NewLedgerRepository(db *gorm.DB) LedgerRepository {
	return &gormLedgerRepository{
		db: db,
	}
}