package ledger

import (
	"account/pkg/accountinfo/dto"
	"account/pkg/ledger/model"
	"account/pkg/repository"
	"context"
	"fmt"
)

type Service interface {
	CreateLedgerEntry(ctx context.Context, info *model.Ledger) error
	GetEntries(ctx context.Context, query *dto.LogQuery) ([]*model.Ledger, error)
	//GetAccountsFor(ctx context.Context, accountQuery *dto.AccountQuery) ([]*dto.AccountResponse, error)
}

type ledgerService struct {
	repository repository.LedgerRepository
}

func (ls *ledgerService) CreateLedgerEntry(ctx context.Context, info *model.Ledger) error {
	if info.Activity == dto.Credit {
		err := ls.repository.CreateLedgerEntry(ctx, info)
		if err != nil {
			return fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
		}
	} else if info.Activity == dto.Debit   {

	}


	return nil
}

func (ls *ledgerService) GetEntries(ctx context.Context, query *dto.LogQuery) ([]*model.Ledger, error) {
	entries, err := ls.repository.GetEntries(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
	}

	return entries, nil
}

//func (sis *ledgerService) GetAccountsFor(ctx context.Context, accountQuery *dto.AccountQuery) ([]*dto.AccountResponse, error) {
//	accountInfos, err := sis.repository.GetAccountData(ctx, accountQuery)
//	if err != nil {
//		return nil, fmt.Errorf("Service.GetAccountsFor", err)
//	}
//
//	return mapper.GetFormattedResponseFor(accountInfos), nil
//}

func NewLedgerService(repository repository.LedgerRepository) Service {
	return &ledgerService{
		repository: repository,
	}
}
