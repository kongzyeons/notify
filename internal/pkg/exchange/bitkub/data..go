package bitkub

type BuyReq struct {
	Market string  `json:"market"`
	Sym    string  `json:"sym"`
	Amount float64 `json:"amount"`
}
type BuyRes struct {
	Error  int `json:"error"`
	Result struct {
		ID            string  `json:"id"`   // Order ID
		Hash          string  `json:"hash"` // Order hash
		Type          string  `json:"typ"`  // Order type (e.g., "limit")
		Amount        float64 `json:"amt"`  // Spending amount
		Rate          float64 `json:"rat"`  // Rate
		Fee           float64 `json:"fee"`  // Fee
		Credit        float64 `json:"cre"`  // Fee credit used
		ReceiveAmount float64 `json:"rec"`  // Amount to receive
		Timestamp     string  `json:"ts"`   // Timestamp
		ClientID      string  `json:"ci"`   // Input client ID

	} `json:"result"`
}

type SellReq struct {
	Market string  `json:"market"`
	Sym    string  `json:"sym"`
	Amount float64 `json:"amount"`
}
type SellRes struct {
	Error  int `json:"error"`
	Result struct {
		ID            string  `json:"id"`   // Order ID
		Hash          string  `json:"hash"` // Order hash
		Type          string  `json:"typ"`  // Order type (e.g., "limit")
		Amount        float64 `json:"amt"`  // Spending amount
		Rate          float64 `json:"rat"`  // Rate
		Fee           float64 `json:"fee"`  // Fee
		Credit        float64 `json:"cre"`  // Fee credit used
		ReceiveAmount float64 `json:"rec"`  // Amount to receive
		Timestamp     string  `json:"ts"`   // Timestamp
		ClientID      string  `json:"ci"`   // Input client ID

	} `json:"result"`
}

type GetsymbolsRes struct {
	Error  int `json:"error"`
	Result []struct {
		ID     int    `json:"id"`
		Symbol string `json:"symbol"`
		Info   string `json:"info"`
	} `json:"result"`
}

type GetMyOpenOrderRes struct {
	Error  int `json:"error"`
	Result []struct {
		ID        string  `json:"id"`
		Hash      string  `json:"hash"`
		Side      string  `json:"side"`
		Type      string  `json:"type"`
		Rate      float64 `json:"rate"`
		Fee       float64 `json:"fee"`
		Credit    float64 `json:"credit"`
		Amount    float64 `json:"amount"`
		Receive   float64 `json:"receive"`
		ParentID  int     `json:"parent_id"`
		SuperID   int     `json:"super_id"`
		ClientID  string  `json:"client_id"`
		Timestamp int64   `json:"ts"`
	} `json:"result"`
}

type GetFiatAccountRes struct {
	Error  int `json:"error"`
	Result []struct {
		ID   string `json:"id"`
		Bank string `json:"bank"`
		Name string `json:"name"`
		Time int64  `json:"time"`
	} `json:"result"`
	Pagination struct {
		Page string `json:"page"`
		Last int64  `json:"last"`
	} `json:"pagination"`
}

type WithdrawReq struct {
	ID  string  `json:"id"`
	Amt float64 `json:"amt"`
}
type WithdrawRes struct {
	Error  int `json:"error"`
	Result struct {
		Txn string  `json:"txn"` // local transaction id
		Acc string  `json:"acc"` // bank account id
		Cur string  `json:"cur"` // currency
		Amt string  `json:"amt"` // withdraw amount
		Fee float64 `json:"fee"` // withdraw fee
		Rec float64 `json:"rec"` // amount to receive
		Ts  float64 `json:"ts"`  // timestamp
	} `json:"result"`
}

type GetBalanceRes struct {
	Error  int `json:"error"`
	Result map[string]struct {
		Available float64 `json:"available"`
		Reserved  float64 `json:"reserved"`
	} `json:"result"`
}

type GetWalletRes struct {
	Error  int                `json:"error"`
	Result map[string]float64 `json:"result"`
}

type Ticker struct {
	ID            int     `json:"id"`
	Last          float64 `json:"last"`
	LowestAsk     float64 `json:"lowestAsk"`
	HighestBid    float64 `json:"highestBid"`
	PercentChange float64 `json:"percentChange"`
	BaseVolume    float64 `json:"baseVolume"`
	QuoteVolume   float64 `json:"quoteVolume"`
	IsFrozen      int     `json:"isFrozen"`
	High24Hr      float64 `json:"high24hr"`
	Low24Hr       float64 `json:"low24hr"`
}

type GetTradingviewHisReq struct {
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Limit     int64  `json:"limit"`
}

type GetTradingviewHisRes struct {
	Success  string    `jsom:"s"`
	Close    []float64 `json:"c"`
	High     []float64 `json:"h"`
	Low      []float64 `json:"l"`
	Open     []float64 `json:"o"`
	Datetime []int64   `json:"t"`
	Volume   []float64 `json:"v"`
}
