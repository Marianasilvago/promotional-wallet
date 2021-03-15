package config

import (
	"fmt"
	"strings"
)

const (
	Delimeter = ","
)

type DataRefresherConfig struct {
	tickerIntervalInSec int
	fsyms               []string
	tsyms               []string
}

func newDataRefresherConfig() DataRefresherConfig {
	fsyms := strings.Split(getString("REFRESHER_FSYMS", ""), Delimeter)
	tsyms := strings.Split(getString("REFRESHER_TSYMS", ""), Delimeter)
	fmt.Println(fsyms)
	fmt.Println(tsyms)
	return DataRefresherConfig{
		tickerIntervalInSec: getInt("TICKER_INTERVAL_IN_SEC"),
		fsyms:               fsyms,
		tsyms:               tsyms,
	}
}

func (cc DataRefresherConfig) GetTickerIntervalInSec() int {
	return cc.tickerIntervalInSec
}

func (cc DataRefresherConfig) GetFsyms() []string {
	return cc.fsyms
}

func (cc DataRefresherConfig) GetTsyms() []string {
	return cc.tsyms
}
