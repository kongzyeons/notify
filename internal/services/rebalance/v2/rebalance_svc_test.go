package rebalance_svc

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func InitConfigTest() error {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	return err
}

func TestRun(t *testing.T) {
	InitConfigTest()
	NewRebalaceSvc(Request{
		Name:             "btc_thb",
		Asset:            "BTC",
		GetTicker:        "THB_BTC",
		PathName:         "../../data/",
		FileName:         "portf_btc_thb.csv",
		Sym:              "btc_thb",
		Timeframe:        "5",
		RebalanceRatio:   float64(50) / float64(100),
		RebalancePercent: 1,
		PeddingRatio:     1,
		RatioWithdraw:    1.5,
	},
	).Run()
}
