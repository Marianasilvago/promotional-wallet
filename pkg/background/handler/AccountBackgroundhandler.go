package handler

import (
	"account/pkg/accountinfo"
	"account/pkg/config"
	"account/pkg/ledger"
	"context"
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

	err := abh.ledgerSvc.ExpireCredits(ctx)
	if err != nil {
		return err
	}
	abh.lgr.Debug("msg", zap.String("eventCode", "ACCOUNT_UPDATED"))
	return nil
}


