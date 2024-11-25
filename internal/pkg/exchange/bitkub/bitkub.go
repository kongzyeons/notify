package bitkub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go_notify/internal/config"
	"go_notify/internal/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BitkubAPI interface {
	Buy(req BuyReq) (res BuyRes, err error)
	Sell(req SellReq) (res SellRes, err error)
	GetWallet() (res GetWalletRes, err error)
	GetTicker(sym ...string) (res map[string]Ticker, err error)
	GetFiatAccount(page, lmt int64) (res GetFiatAccountRes, err error)
	Withdraw(req WithdrawReq) (res WithdrawRes, err error)

	// optional
	GetTradingviewHis(req GetTradingviewHisReq) (res GetTradingviewHisRes, err error)
	GetListSymbols() (res GetsymbolsRes, err error)
	GetTimeServer() (res int64, err error)
	GetMyOpenOrder(sym string) (res GetMyOpenOrderRes, err error)
	GetBalance() (res GetBalanceRes, err error)
}

type bitkubAPI struct {
	apiUrl          string
	BitkubAPIKey    string
	BitkubAPISecret string
}

func NewBitkubAPI() BitkubAPI {
	cfg := config.InitConfig()
	return &bitkubAPI{
		apiUrl:          "https://api.bitkub.com",
		BitkubAPIKey:    cfg.BitkubAPIKey,
		BitkubAPISecret: cfg.BitkubAPISecret,
	}
}

func (self *bitkubAPI) Buy(req BuyReq) (res BuyRes, err error) {
	path := "/api/v3/market/place-bid"
	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	param := map[string]interface{}{
		"sym": req.Sym,
		"amt": req.Amount,
		"rat": 0.0,
		"typ": req.Market,
	}

	payload := []string{
		ts,
		http.MethodPost,
		path,
		string(mustMarshal(param)),
	}
	sig := genSign(self.BitkubAPISecret, strings.Join(payload, ""))

	request, err := http.NewRequest(http.MethodPost, self.apiUrl+path, bytes.NewReader(mustMarshal(param)))
	// request, err := http.NewRequest(http.MethodPost, self.apiUrl+path+queryParam, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BTK-TIMESTAMP", ts)
	request.Header.Set("X-BTK-SIGN", sig)
	request.Header.Set("X-BTK-APIKEY", self.BitkubAPIKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}

	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}

func (self *bitkubAPI) Sell(req SellReq) (res SellRes, err error) {
	path := "/api/v3/market/place-ask"
	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	param := map[string]interface{}{
		"sym": req.Sym,
		"amt": req.Amount,
		"rat": 0.0,
		"typ": req.Market,
	}

	payload := []string{
		ts,
		http.MethodPost,
		path,
		string(mustMarshal(param)),
	}
	sig := genSign(self.BitkubAPISecret, strings.Join(payload, ""))

	request, err := http.NewRequest(http.MethodPost, self.apiUrl+path, bytes.NewReader(mustMarshal(param)))
	// request, err := http.NewRequest(http.MethodPost, self.apiUrl+path+queryParam, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BTK-TIMESTAMP", ts)
	request.Header.Set("X-BTK-SIGN", sig)
	request.Header.Set("X-BTK-APIKEY", self.BitkubAPIKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}

	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}

func (self *bitkubAPI) GetMyOpenOrder(sym string) (res GetMyOpenOrderRes, err error) {
	path := "/api/v3/market/my-open-orders"

	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	param := map[string]interface{}{
		"sym": sym,
	}

	queryParam := genQueryParam(self.apiUrl+path, param)
	payload := []string{ts, "GET", path, queryParam}
	sig := genSign(self.BitkubAPISecret, strings.Join(payload, ""))

	req, err := http.NewRequest("GET", self.apiUrl+path+queryParam, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BTK-TIMESTAMP", ts)
	req.Header.Set("X-BTK-SIGN", sig)
	req.Header.Set("X-BTK-APIKEY", self.BitkubAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct

	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}

func (self *bitkubAPI) GetFiatAccount(page, lmt int64) (res GetFiatAccountRes, err error) {
	path := "/api/v3/fiat/accounts/"

	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	param := map[string]interface{}{
		"p":   page,
		"lmt": lmt,
	}

	queryParam := genQueryParam(self.apiUrl+path, param)
	payload := []string{ts, http.MethodPost, path, queryParam}
	sig := genSign(self.BitkubAPISecret, strings.Join(payload, ""))

	request, err := http.NewRequest(http.MethodPost, self.apiUrl+path+queryParam, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BTK-TIMESTAMP", ts)
	request.Header.Set("X-BTK-SIGN", sig)
	request.Header.Set("X-BTK-APIKEY", self.BitkubAPIKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}

func (self *bitkubAPI) Withdraw(req WithdrawReq) (res WithdrawRes, err error) {
	path := "/api/v3/fiat/withdraw"
	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	param := map[string]interface{}{
		"id":  req.ID,
		"amt": req.Amt,
	}

	payload := []string{
		ts,
		http.MethodPost,
		path,
		string(mustMarshal(param)),
	}
	sig := genSign(self.BitkubAPISecret, strings.Join(payload, ""))

	request, err := http.NewRequest(http.MethodPost, self.apiUrl+path, bytes.NewReader(mustMarshal(param)))
	if err != nil {
		log.Println(err)
		return res, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BTK-TIMESTAMP", ts)
	request.Header.Set("X-BTK-SIGN", sig)
	request.Header.Set("X-BTK-APIKEY", self.BitkubAPIKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}

	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}

func (self *bitkubAPI) GetBalance() (res GetBalanceRes, err error) {
	path := "/api/v3/market/balances"

	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	param := map[string]interface{}{}

	queryParam := genQueryParam(self.apiUrl+path, param)
	payload := []string{ts, http.MethodPost, path, queryParam}
	sig := genSign(self.BitkubAPISecret, strings.Join(payload, ""))

	request, err := http.NewRequest(http.MethodPost, self.apiUrl+path+queryParam, nil)

	if err != nil {
		log.Println(err)
		return res, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BTK-TIMESTAMP", ts)
	request.Header.Set("X-BTK-SIGN", sig)
	request.Header.Set("X-BTK-APIKEY", self.BitkubAPIKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct

	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	result := make(map[string]struct {
		Available float64 `json:"available"`
		Reserved  float64 `json:"reserved"`
	})
	if len(res.Result) > 0 {
		for key, value := range res.Result {
			if value.Available > 0 || value.Reserved > 0 {
				fmt.Println(fmt.Sprintf("%s : %v", key, value))
				result[key] = value
			}
		}
		res.Result = result
	}
	return res, nil
}

func (self *bitkubAPI) GetWallet() (res GetWalletRes, err error) {
	path := "/api/v3/market/wallet"

	ts := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	param := map[string]interface{}{}

	queryParam := genQueryParam(self.apiUrl+path, param)
	payload := []string{ts, http.MethodPost, path, queryParam}
	sig := genSign(self.BitkubAPISecret, strings.Join(payload, ""))

	request, err := http.NewRequest(http.MethodPost, self.apiUrl+path+queryParam, nil)

	if err != nil {
		log.Println(err)
		return res, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BTK-TIMESTAMP", ts)
	request.Header.Set("X-BTK-SIGN", sig)
	request.Header.Set("X-BTK-APIKEY", self.BitkubAPIKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct

	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	result := make(map[string]float64)
	if len(res.Result) > 0 {
		for k, v := range res.Result {
			if v > 0 {
				result[k] = v
			}
		}
		res.Result = result
	}

	return res, nil
}

func (self *bitkubAPI) GetTicker(sym ...string) (res map[string]Ticker, err error) {
	path := "/api//market/ticker"
	param := map[string]interface{}{}
	if len(sym) > 0 {
		if sym[0] != "" {
			param["sym"] = sym[0]
		}
	}
	queryParam := genQueryParam(self.apiUrl+path, param)

	request, err := http.NewRequest(http.MethodGet, self.apiUrl+path+queryParam, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}

func (self *bitkubAPI) GetTradingviewHis(req GetTradingviewHisReq) (res GetTradingviewHisRes, err error) {
	path := "/tradingview/history"

	if !utils.ValueInSlice[string](req.Timeframe, []string{"1", "5", "15", "60", "240", "1D"}) {
		msg := "error validate timeframe"
		log.Println(msg)
		return res, errors.New(msg)
	}

	resolutionSeconds, _ := convertCustomToSeconds(req.Timeframe)
	toTimestamp := int64(time.Now().UTC().Unix())

	param := map[string]interface{}{
		"symbol":     req.Symbol,
		"resolution": req.Timeframe,
		"from":       toTimestamp - ((req.Limit - 1) * int64(resolutionSeconds)),
		"to":         toTimestamp,
	}

	queryParam := genQueryParam(self.apiUrl+path, param)

	request, err := http.NewRequest(http.MethodGet, self.apiUrl+path+queryParam, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}

func (self *bitkubAPI) GetListSymbols() (res GetsymbolsRes, err error) {
	apiUrl := fmt.Sprintf("%s%s", self.apiUrl, "/api/market/symbols")
	request, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (self *bitkubAPI) GetTimeServer() (res int64, err error) {
	apiUrl := fmt.Sprintf("%s%s", self.apiUrl, "/api/v3/servertime")
	request, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("api bitkub error status : %d", resp.StatusCode))
		return res, err
	}

	// Read the response body
	bodyRes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return res, err
	}
	// Unmarshal the JSON response into the struct
	if err := json.Unmarshal(bodyRes, &res); err != nil {
		log.Println(err)
		return res, err
	}

	return res, err
}
