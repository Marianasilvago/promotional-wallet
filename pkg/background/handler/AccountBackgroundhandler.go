package handler

import (
	"account/pkg/accountinfo"
	"account/pkg/accountinfo/model"
	"account/pkg/config"
	"account/pkg/ledger"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type AccountBackgroundHandler struct {
	lgr        *zap.Logger
	accountSvc accountinfo.Service
	ledgerSvc  ledger.Service
	drc        config.DataRefresherConfig
}

func NewAccountBackgroundHandler(lgr *zap.Logger, accountSvc accountinfo.Service, ledgerSvc ledger.Service, drc config.DataRefresherConfig) *AccountBackgroundHandler {
	return &AccountBackgroundHandler{
		lgr:        lgr,
		accountSvc: accountSvc,
		ledgerSvc:  ledgerSvc,
		drc:        drc,
	}
}

func (abh *AccountBackgroundHandler) UpdateAccountExpiredEntries() error {
	ctx := context.Background()
	abh.lgr.Sugar().Infof("Updating account data..")

	expiredEntries, err := abh.ledgerSvc.ExpireCredits(ctx)
	if err != nil {
		return err
	}
	for _, expiredEntry := range expiredEntries {
		fmt.Println(fmt.Sprintf("=================%+v", expiredEntry))
		accountInfo := &model.AccountInfo{
			ID:      expiredEntry.AccountID,
			Balance: expiredEntry.Amount * -1,
		}
		err := abh.accountSvc.CreateOrUpdateAccountInfo(ctx, accountInfo)
		if err != nil {
			return fmt.Errorf("AccountBackgroundHandler.UpdateAccountExpiredEntries . error %v", err)
		}
	}

	abh.lgr.Debug("msg", zap.String("eventCode", "ACCOUNT_UPDATED"))
	return nil
}
