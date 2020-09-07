package test

import (
	"portfolio"
	"time"
)

var (
	testTenant     = "dv2"
	testPortfolio1 = portfolio.Portfolio{
		TenantId:         testTenant,
		Id:               "1",
		Name:             "An ETF Stocks portfolio",
		Description:      "Desc1",
		Assets:           []portfolio.Asset{{"SPY", 33}, {"QQQ", 33}, {"VFV", 34}},
		CreationDate:     time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
		LastModifiedDate: time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
	}
	testPortfolio2 = portfolio.Portfolio{
		TenantId:         testTenant,
		Id:               "2",
		Name:             "An ETF Bonds portfolio",
		Description:      "Desc2",
		Assets:           []portfolio.Asset{{"IEF", 50}, {"SHY", 50}},
		CreationDate:     time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
		LastModifiedDate: time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
	}
	testPortfolios = []portfolio.Portfolio{testPortfolio1, testPortfolio2}
)
