package ema_svc

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func InitConfigTest() error {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	return err
}

func TestRun(t *testing.T) {
	InitConfigTest()

	NewEmaSvc(
		Request{
			Name:      "btc_thb",
			Asset:     "BTC",
			GetTicker: "THB_BTC",
			Sym:       "btc_thb",
			Timeframe: "15"},
	).Run()
}
