package main

import (
	"go_notify/internal/pkg/conjob"
	ema_svc "go_notify/internal/services/ema"
	rebalance_svc "go_notify/internal/services/rebalance/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func createScheduler() {
	log.Println("start scheduler...")
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Println(err)
		return
	}

	conjob.NewJob(
		s, "bitkub rebalance btc every 5 minitute",
		5*time.Minute,
		rebalance_svc.NewRebalaceSvc(
			rebalance_svc.Request{
				Name:             "btc_thb",
				Asset:            "BTC",
				GetTicker:        "THB_BTC",
				PathName:         "internal/services/data/",
				FileName:         "portf_btc_thb.csv",
				Sym:              "btc_thb",
				Timeframe:        "5",
				RebalanceRatio:   float64(50) / float64(100),
				RebalancePercent: 1,
				PeddingRatio:     1,
				RatioWithdraw:    1.5,
			},
		),
	)

	conjob.NewJob(
		s, "bitkub ema signal btc every 15 minitute",
		15*time.Minute,
		ema_svc.NewEmaSvc(
			ema_svc.Request{
				Name:      "btc_thb",
				Sym:       "btc_thb",
				Timeframe: "15",
			},
		),
	)

	// Start the scheduler
	s.Start()
	// Set up a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("\nInterrupt signal received. Exiting...")
		_ = s.Shutdown()
		os.Exit(0)
	}()
	for {

	}
}
