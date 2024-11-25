package ema_svc

import (
	"encoding/json"
	"fmt"
	"go_notify/internal/models"
	"go_notify/internal/pkg/exchange/bitkub"
	"go_notify/internal/pkg/line"
	"log"
	"time"
)

type EmaSvc interface {
	Run() error
}

type emaSvc struct {
	lineAPI   line.LineAPI
	bitkubAPI bitkub.BitkubAPI
	loc       *time.Location
	params    Request
}

func NewEmaSvc(params Request) EmaSvc {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Println(err)
		loc = time.UTC
	}
	return &emaSvc{
		lineAPI:   line.NewLineAPI(),
		bitkubAPI: bitkub.NewBitkubAPI(),
		loc:       loc,
		params:    params,
	}
}

func (self *emaSvc) Run() error {

	df, err := self.bitkubAPI.GetTradingviewHis(bitkub.GetTradingviewHisReq{
		Symbol:    self.params.Sym,
		Timeframe: self.params.Timeframe,
		Limit:     999,
	})
	if err != nil {
		log.Println(err)
		res := models.Response[any]{
			Status:  false,
			Code:    500,
			Message: "error bitkub api get view history",
			Data:    nil,
		}
		jsonData, _ := json.Marshal(res)
		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
			log.Println(err)
			return err
		}
		return err
	}
	count := len(df.Close)
	ema13 := ema(df.Close, 13)
	emaFastA := ema13[count-2]
	emaFastB := ema13[count-3]
	// fmt.Println("EMA_fast_A = ", emaFastA)
	// fmt.Println("EMA_fast_A = ", emaFastB)

	ema33 := ema(df.Close, 33)
	emaSlowA := ema33[count-2]
	emaSlowB := ema33[count-3]
	// fmt.Println("EMA_slow_A = ", emaSlowA)
	// fmt.Println("EMA_slow_B = ", emaSlowB)

	trend := NonTrend
	signal := NoneSignal
	if emaFastA > emaSlowA {
		trend = UpTrend
		if emaFastB < emaSlowB {
			signal = BuySignal
		}
	} else if emaFastA < emaSlowA {
		trend = DownTrend
		if emaFastB > emaSlowB {
			signal = SellSignal
		}
	}

	if trend != NonTrend && signal != NoneSignal {
		res := models.Response[Response]{
			Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
			Status:  true,
			Code:    200,
			Message: "signal",
			Data: &Response{
				DateTime: time.Now().In(self.loc),
				Name:     self.params.Name,
				Signal:   signal,
				Trend:    trend,
			},
		}
		jsonData, _ := json.MarshalIndent(res, "", "   ")

		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
