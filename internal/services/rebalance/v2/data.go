package rebalance_svc

import "time"

// type Response struct {
// 	THB   float64 `json:"thb"`
// 	BTC   float64 `json:"btc"`
// 	Value float64 `json:"value"`
// }

type Request struct {
	Name             string  `json:"name"`
	Asset            string  `json:"asset"`
	GetTicker        string  `json:"getTicker"`
	Sym              string  `json:"sym"`
	PathName         string  `json:"pathName"`
	FileName         string  `json:"fileName"`
	Timeframe        string  `json:"timeframe"`
	RebalanceRatio   float64 `json:"rebalanceRatio"`
	RebalancePercent float64 `json:"rebalancePercent"`
	PeddingRatio     float64 `json:"peddingRatio"`
	RatioWithdraw    float64 `json:"ratioWithdraw"`
}

type ResponseSignal struct {
	DateTime time.Time `json:"dateTime"`
	Name     string    `json:"name"`
	Signal   Signal    `json:"signal"`
	Trend    Signal    `json:"trend"`
}

type ResponseReblance struct {
	DateTime      time.Time `json:"dateTime"`
	AssetPrice    float64   `json:"assetPrice"`
	Units         float64   `json:"units"`
	AssetValue    float64   `json:"assetValue"`
	Cash          float64   `json:"cash"`
	Total         float64   `json:"total"`
	RebalanceMark float64   `json:"rebalanceMark"`
	Status        Status    `json:"status"`
	Diff          float64   `json:"diff"`
	CashFlow      float64   `json:"cashFlow"`
	NewAssetValue float64   `json:"newAssetValue"`
	NewUnits      float64   `json:"newUnits"`
	NewCash       float64   `json:"newCash"`
	NewTotal      float64   `json:"newTotal"`
}

type ResponseWithdraw struct {
	DateTime time.Time `json:"dateTime"`
	Name     string    `json:"name"`
	Bank     string    `json:"bank"`
	CashFlow float64   `json:"cashFlow"`
	Fee      float64   `json:"fee"`     // withdraw fee
	Receive  float64   `json:"receive"` // amount to receive
}

type Status string

const (
	Buy       Status = "buy"
	Sell      Status = "sell"
	NoneTrade Status = "non trade"
	Withdraw  Status = "withdraw"
)

type Signal string

const (
	BuySignal  Signal = "buySignal"
	UpTrend    Signal = "upTrend"
	SellSignal Signal = "sellSignal"
	DownTrend  Signal = "downTrend"
	NoneSignal Signal = "none"
	NonTrend   Signal = "none"
)
