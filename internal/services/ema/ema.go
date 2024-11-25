package ema_svc

import "time"

type Request struct {
	Name      string `json:"name"`
	Asset     string `json:"asset"`
	GetTicker string `json:"getTicker"`
	Sym       string `json:"sym"`
	FileName  string `json:"fileName"`
	Timeframe string `json:"timeframe"`
}

type Response struct {
	DateTime time.Time `json:"dateTime"`
	Name     string    `json:"name"`
	Signal   Signal    `json:"signal"`
	Trend    Signal    `json:"trend"`
}

type Signal string

const (
	BuySignal  Signal = "buySignal"
	UpTrend    Signal = "upTrend"
	SellSignal Signal = "sellSignal"
	DownTrend  Signal = "downTrend"
	NoneSignal Signal = "none"
	NonTrend   Signal = "none"
)
