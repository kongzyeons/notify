package rebalance_svc

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"go_notify/internal/models"
	"go_notify/internal/pkg/exchange/bitkub"
	"go_notify/internal/pkg/line"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type RebalaceSvc interface {
	Run() error
}

type rebalaceSvc struct {
	lineAPI   line.LineAPI
	bitkubAPI bitkub.BitkubAPI
	loc       *time.Location
	pathName  string
	params    Request
}

func NewRebalaceSvc(params Request) RebalaceSvc {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Println(err)
		loc = time.UTC
	}

	return &rebalaceSvc{
		lineAPI:   line.NewLineAPI(),
		bitkubAPI: bitkub.NewBitkubAPI(),
		loc:       loc,
		pathName:  params.PathName,
		// pathName: "../../data/",
		params: params,
	}
}

func (self *rebalaceSvc) Run() error {
	wallet, err := self.bitkubAPI.GetWallet()
	if err != nil {
		log.Println(err)
		res := models.Response[any]{
			Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
			Status:  false,
			Code:    500,
			Message: "error bitkub api get wallet",
			Data:    nil,
		}
		jsonData, _ := json.MarshalIndent(res, "", "   ")
		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
			log.Println(err)
			return err
		}
		return err
	}
	if len(wallet.Result) < 0 {
		msg := "not found data from bitkub"
		log.Println(msg)
		res := models.Response[any]{
			Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
			Status:  false,
			Code:    500,
			Message: msg,
			Data:    nil,
		}
		jsonData, _ := json.MarshalIndent(res, "", "   ")
		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
			log.Println(err)
			return err
		}
		return errors.New(msg)
	}

	var thb float64
	if _, ok := wallet.Result["THB"]; ok {
		thb = wallet.Result["THB"]
	}
	var asset float64
	if _, ok := wallet.Result[self.params.Asset]; ok {
		asset = wallet.Result[self.params.Asset]
	}

	var total float64
	var assetPrice float64
	for key, value := range wallet.Result {
		if key == "THB" {
			total += value
		} else {
			price, err := self.bitkubAPI.GetTicker(fmt.Sprintf("THB_%s", key))
			if err != nil {
				log.Println(err)
				res := models.Response[any]{
					Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
					Status:  false,
					Code:    500,
					Message: "error bitkub api get ticker",
					Data:    nil,
				}
				jsonData, _ := json.MarshalIndent(res, "", "   ")
				if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
					log.Println(err)
					return err
				}
				return err
			}
			if fmt.Sprintf("THB_%s", key) == self.params.GetTicker {
				assetPrice = price[self.params.GetTicker].Last
			}
			total += value * price[fmt.Sprintf("THB_%s", key)].Last
		}
	}

	// records, err := self.readCSV(self.pathName + self.params.FileName)
	// if err != nil {
	// 	log.Println(err)
	// 	res := models.Response[any]{
	// 		Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
	// 		Status:  false,
	// 		Code:    500,
	// 		Message: "error read csv",
	// 		Data:    nil,
	// 	}
	// 	jsonData, _ := json.MarshalIndent(res, "", "   ")
	// 	if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	return err
	// }
	// if len(records) == 0 {
	// 	log.Println("not found data")
	// 	res := models.Response[any]{
	// 		Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
	// 		Status:  false,
	// 		Code:    500,
	// 		Message: "not found data",
	// 		Data:    nil,
	// 	}
	// 	jsonData, _ := json.MarshalIndent(res, "", "   ")
	// 	if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	return err
	// }

	var data ResponseReblance
	data.DateTime = time.Now().In(self.loc)
	data.AssetPrice = assetPrice
	data.Units = asset
	data.AssetValue = asset * assetPrice
	data.Cash = thb
	data.Total = total
	data.RebalanceMark = data.Total * self.params.RebalanceRatio

	// fmt.Println(data.AssetValue, (data.RebalanceMark))
	// fmt.Println(data.AssetValue, (data.RebalanceMark + (data.RebalanceMark * self.params.RebalancePercent / 100)))
	// fmt.Println(data.AssetValue, (data.RebalanceMark - (data.RebalanceMark * self.params.RebalancePercent / 100)))

	// Asset_01_Value > (Rebalance_mark01 + (Rebalance_mark01 *Rebalance_percent/100) ) :
	if data.AssetValue > (data.RebalanceMark + (data.RebalanceMark * self.params.RebalancePercent / 100)) {
		data.Status = NoneTrade
		data.Diff = data.AssetValue - data.RebalanceMark
		data.CashFlow = data.Diff / 2
		data.NewAssetValue = data.AssetValue - data.Diff
		data.NewUnits = data.NewAssetValue / data.AssetPrice
		data.NewCash = data.Cash + data.Diff - data.CashFlow
		data.NewTotal = data.NewAssetValue + data.NewCash
		if data.Diff > 10*self.params.PeddingRatio {
			data.Status = Sell
			res, err := self.bitkubAPI.Sell(bitkub.SellReq{
				Market: "market",
				Sym:    self.params.Sym,
				Amount: data.Diff / assetPrice,
			})
			if (err != nil) || res.Error != 0 {
				msg := fmt.Sprintf("%v %d", err, res.Error)
				log.Println(fmt.Sprintf("error bitkub api sell : %s", msg))
				res := models.Response[any]{
					Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
					Status:  false,
					Code:    500,
					Message: fmt.Sprintf("error bitkub api sell : %s", msg),
					Data:    nil,
				}
				jsonData, _ := json.MarshalIndent(res, "", "   ")
				if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
					log.Println(err)
					return err
				}
				return nil
			}
		}
		// Asset_01_Value < (Rebalance_mark01 - (Rebalance_mark01*Rebalance_percent/100) ) :
	} else if data.AssetValue < (data.RebalanceMark - (data.RebalanceMark * self.params.RebalancePercent / 100)) {
		data.Status = NoneTrade
		data.Diff = data.RebalanceMark - data.AssetValue
		data.CashFlow = 0
		data.NewAssetValue = data.AssetValue + data.Diff
		data.NewUnits = data.NewAssetValue / data.AssetPrice
		data.NewCash = data.Cash - data.Diff
		data.NewTotal = data.NewAssetValue + data.NewCash
		if data.Diff > 10*self.params.PeddingRatio {
			data.Status = Buy
			res, err := self.bitkubAPI.Buy(bitkub.BuyReq{
				Market: "market",
				Sym:    self.params.Sym,
				Amount: data.Diff,
			})
			if (err != nil) || res.Error != 0 {
				msg := fmt.Sprintf("%v %d", err, res.Error)
				log.Println(fmt.Sprintf("error bitkub api buy : %s", msg))
				res := models.Response[any]{
					Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
					Status:  false,
					Code:    500,
					Message: fmt.Sprintf("error bitkub api buy : %s", msg),
					Data:    nil,
				}
				jsonData, _ := json.MarshalIndent(res, "", "   ")
				if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
					log.Println(err)
					return err
				}
				return nil
			}
		}

	} else {
		data.Status = NoneTrade
	}

	if data.Status != NoneTrade {
		res := models.Response[ResponseReblance]{
			Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
			Status:  true,
			Code:    200,
			Message: string(data.Status),
			Data:    &data,
		}
		jsonData, _ := json.MarshalIndent(res, "", "   ")
		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
			log.Println(err)
			return err
		}
		// if err := self.writeCSV(self.pathName+self.params.FileName, data, records); err != nil {
		// 	log.Println(err)
		// 	return err
		// }
	}

	//withdraw
	if data.Status == Sell {
		if data.CashFlow > 20*self.params.RatioWithdraw {
			cashFlow := math.Round(data.CashFlow*100) / 100
			account, err := self.bitkubAPI.GetFiatAccount(1, 1)
			if err != nil || account.Error != 0 {
				msg := fmt.Sprintf("%v %d", err, account.Error)
				log.Println(msg)
				res := models.Response[any]{
					Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
					Status:  false,
					Code:    500,
					Message: fmt.Sprintf("error bitkub api get fiat account : %s", msg),
					Data:    nil,
				}
				jsonData, _ := json.MarshalIndent(res, "", "   ")
				if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
					log.Println(err)
					return err
				}
				return err
			}

			resWithdraw, err := self.bitkubAPI.Withdraw(bitkub.WithdrawReq{
				ID:  account.Result[0].ID,
				Amt: cashFlow,
			})
			if err != nil || resWithdraw.Error != 0 {
				msg := fmt.Sprintf("%v %d", err, resWithdraw.Error)
				res := models.Response[any]{
					Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
					Status:  false,
					Code:    500,
					Message: fmt.Sprintf("error bitkub api withdraw : %s", msg),
					Data:    nil,
				}
				jsonData, _ := json.MarshalIndent(res, "", "   ")
				if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
					log.Println(err)
					return err
				}
			}

			// success
			data := ResponseWithdraw{
				DateTime: time.Now().In(self.loc),
				Name:     account.Result[0].Name,
				Bank:     account.Result[0].Bank,
				CashFlow: cashFlow,
				Fee:      resWithdraw.Result.Fee,
				Receive:  resWithdraw.Result.Rec,
			}
			res := models.Response[ResponseWithdraw]{
				Title:   fmt.Sprintf("Rebalance Bitkub %s", self.params.Name),
				Status:  true,
				Code:    200,
				Message: string(Withdraw),
				Data:    &data,
			}
			jsonData, _ := json.MarshalIndent(res, "", "   ")
			if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (self *rebalaceSvc) checkSignal() (trend, signal Signal, err error) {
	df, err := self.bitkubAPI.GetTradingviewHis(bitkub.GetTradingviewHisReq{
		Symbol:    self.params.Sym,
		Timeframe: self.params.Timeframe,
		Limit:     999,
	})
	if err != nil {
		log.Println(err)
		return NonTrend, NoneSignal, err
	}

	count := len(df.Close)

	ema13 := ema(df.Close, 13)
	emaFastA := ema13[count-2]
	emaFastB := ema13[count-3]

	ema33 := ema(df.Close, 33)
	emaSlowA := ema33[count-2]
	emaSlowB := ema33[count-3]

	trend = NonTrend
	signal = NoneSignal
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
	return trend, signal, nil
}

func (self *rebalaceSvc) readCSV(filename string) (records [][]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return records, err
			}
		} else {
			fmt.Println("Error opening file:", err)
			return records, err
		}
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err = reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return records, err
	}

	if len(records) == 0 {
		header := []string{
			"DateTime", "AssetPrice", "Units", "AssetValue", "Cash", "Total",
			"RebalanceMark", "Status", "Diff",
			"CashFlow",
			"NewAssetValue", "NewUnits", "NewCash", "NewTotal",
		}
		records = append(records, header)
		file.Close()

		// Open the CSV file for writing (truncate the file)
		file, err = os.Create(filename)
		if err != nil {
			fmt.Println("Error opening file for writing:", err)
			return records, err
		}
		defer file.Close()

		// Create a CSV writer
		writer := csv.NewWriter(file)

		// Write all records (including the new row) back to the CSV file
		err = writer.WriteAll(records)
		if err != nil {
			fmt.Println("Error writing to CSV:", err)
			return records, err
		}

		// Flush the writer
		writer.Flush()
	}
	return records, err
}

func (self *rebalaceSvc) writeCSV(filename string, data ResponseReblance, records [][]string) error {
	if len(records) == 0 {
		header := []string{
			"DateTime", "AssetPrice", "Units", "AssetValue", "Cash", "Total",
			"RebalanceMark", "Status", "Diff",
			"CashFlow",
			"NewAssetValue", "NewUnits", "NewCash", "NewTotal",
		}
		records = append(records, header)
	}
	row := []string{
		data.DateTime.In(self.loc).Format("02/01/2006 15:04"),
		strconv.FormatFloat(data.AssetPrice, 'f', 2, 64),
		strconv.FormatFloat(data.Units, 'f', -1, 64),
		strconv.FormatFloat(data.AssetValue, 'f', 2, 64),
		strconv.FormatFloat(data.Cash, 'f', 2, 64),
		strconv.FormatFloat(data.Total, 'f', 2, 64),
		strconv.FormatFloat(data.RebalanceMark, 'f', 2, 64),
		string(data.Status),
		strconv.FormatFloat(data.Diff, 'f', 2, 64),
		strconv.FormatFloat(data.CashFlow, 'f', 2, 64),
		strconv.FormatFloat(data.NewAssetValue, 'f', 2, 64),
		strconv.FormatFloat(data.NewUnits, 'f', -1, 64),
		strconv.FormatFloat(data.NewCash, 'f', 2, 64),
		strconv.FormatFloat(data.NewTotal, 'f', 2, 64),
	}
	records = append(records, row)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write all records (including the new row) back to the CSV file
	err = writer.WriteAll(records)
	if err != nil {
		fmt.Println("Error writing to CSV:", err)
		return err
	}

	// Flush the writer
	writer.Flush()
	return nil
}

// func (self *rebalaceSvc) Run() error {
// 	wallet, err := self.bitkubAPI.GetWallet()
// 	if err != nil {
// 		log.Println(err)
// 		res := models.Response[any]{
// 			Title:   "Rebalance Bitkub btc_thb",
// 			Status:  false,
// 			Code:    500,
// 			Message: "error bitkub api get wallet",
// 			Data:    nil,
// 		}
// 		jsonData, _ := json.MarshalIndent(res, "", "   ")
// 		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		return err
// 	}
// 	if len(wallet.Result) < 0 {
// 		msg := "not found data from bitkub"
// 		log.Println(msg)
// 		res := models.Response[any]{
// 			Title:   "Rebalance Bitkub btc_thb",
// 			Status:  false,
// 			Code:    500,
// 			Message: msg,
// 			Data:    nil,
// 		}
// 		jsonData, _ := json.MarshalIndent(res, "", "   ")
// 		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		return errors.New(msg)
// 	}

// 	var thb float64
// 	if _, ok := wallet.Result["THB"]; ok {
// 		thb = wallet.Result["THB"]
// 	}
// 	var btc float64
// 	if _, ok := wallet.Result["BTC"]; ok {
// 		btc = wallet.Result["BTC"]
// 	}

// 	price, err := self.bitkubAPI.GetTicker("thb_btc")
// 	if err != nil {
// 		log.Println(err)
// 		res := models.Response[any]{
// 			Title:   "Rebalance Bitkub btc_thb",
// 			Status:  false,
// 			Code:    500,
// 			Message: "error bitkub api get ticker",
// 			Data:    nil,
// 		}
// 		jsonData, _ := json.MarshalIndent(res, "", "   ")
// 		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		return err
// 	}

// 	records, err := self.readCSV(self.fileName)
// 	if err != nil {
// 		log.Println(err)
// 		res := models.Response[any]{
// 			Title:   "Rebalance Bitkub btc_thb",
// 			Status:  false,
// 			Code:    500,
// 			Message: "error read csv",
// 			Data:    nil,
// 		}
// 		jsonData, _ := json.MarshalIndent(res, "", "   ")
// 		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		return err
// 	}
// 	if len(records) == 0 {
// 		log.Println("not found data")
// 		res := models.Response[any]{
// 			Title:   "Rebalance Bitkub btc_thb",
// 			Status:  false,
// 			Code:    500,
// 			Message: "not found data",
// 			Data:    nil,
// 		}
// 		jsonData, _ := json.MarshalIndent(res, "", "   ")
// 		if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		return err
// 	}

// 	var data Response
// 	data.DateTime = time.Now().In(self.loc)
// 	data.AssetPrice = price["THB_BTC"].Last
// 	data.Units = btc
// 	data.AssetValue = btc * price["THB_BTC"].Last
// 	data.Cash = thb
// 	data.Total = data.Cash + data.AssetValue
// 	data.RebalanceMark = data.Total / 2
// 	if data.AssetValue > data.RebalanceMark {
// 		data.Status = Sell
// 		data.Diff = data.AssetValue - data.RebalanceMark
// 		data.CashFlow = data.Diff / 2
// 		data.NewAssetValue = data.AssetValue - data.Diff
// 		data.NewUnits = data.NewAssetValue / data.AssetPrice
// 		data.NewCash = data.Cash + data.Diff - data.CashFlow
// 		data.NewTotal = data.NewAssetValue + data.NewCash

// 		res, err := self.bitkubAPI.Sell(bitkub.SellReq{
// 			Market: "market",
// 			Sym:    "btc_thb",
// 			Amount: data.Diff / price["THB_BTC"].Last,
// 		})
// 		if (err != nil) || res.Error != 0 {
// 			msg := fmt.Sprintf("%v %d", err, res.Error)
// 			log.Println(fmt.Sprintf("error bitkub api sell : %s", msg))
// 			// res := models.Response[any]{
// 			// 	Title:   "Rebalance Bitkub btc_thb",
// 			// 	Status:  false,
// 			// 	Code:    500,
// 			// 	Message: fmt.Sprintf("error bitkub api sell : %s", msg),
// 			// 	Data:    nil,
// 			// }
// 			// jsonData, _ := json.MarshalIndent(res, "", "   ")
// 			// if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 			// 	log.Println(err)
// 			// 	return err
// 			// }
// 			// data.Status = NoneTrade
// 			// data.CurrentPriceBuy = currentPriceBuy
// 			// data.CurrentDiffBuy = currentDiffBuy
// 			// if err := self.writeCSV(self.fileName, data, records); err != nil {
// 			// 	log.Println(err)
// 			// 	return err
// 			// }
// 			return nil
// 		}

// 	} else if data.AssetValue < data.RebalanceMark {
// 		data.Status = Buy
// 		data.Diff = data.RebalanceMark - data.AssetValue
// 		data.CashFlow = 0
// 		data.NewAssetValue = data.AssetValue + data.Diff
// 		data.NewUnits = data.NewAssetValue / data.AssetPrice
// 		data.NewCash = data.Cash - data.Diff
// 		data.NewTotal = data.NewAssetValue + data.NewCash

// 		res, err := self.bitkubAPI.Buy(bitkub.BuyReq{
// 			Market: "market",
// 			Sym:    "btc_thb",
// 			Amount: data.Diff,
// 		})
// 		if (err != nil) || res.Error != 0 {
// 			msg := fmt.Sprintf("%v %d", err, res.Error)
// 			log.Println(fmt.Sprintf("error bitkub api buy : %s", msg))
// 			// res := models.Response[any]{
// 			// 	Title:   "Rebalance Bitkub btc_thb",
// 			// 	Status:  false,
// 			// 	Code:    500,
// 			// 	Message: fmt.Sprintf("error bitkub api buy : %s", msg),
// 			// 	Data:    nil,
// 			// }
// 			// jsonData, _ := json.MarshalIndent(res, "", "   ")
// 			// if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 			// 	log.Println(err)
// 			// 	return err
// 			// }
// 			// data.Status = NoneTrade
// 			// data.CurrentPriceBuy = currentPriceBuy
// 			// data.CurrentDiffBuy = currentDiffBuy
// 			// if err := self.writeCSV(self.fileName, data, records); err != nil {
// 			// 	log.Println(err)
// 			// 	return err
// 			// }
// 			return nil
// 		}

// 	} else {
// 		data.Status = NoneTrade
// 	}

// 	res := models.Response[Response]{
// 		Title:   "Rebalance Bitkub btc_thb",
// 		Status:  true,
// 		Code:    200,
// 		Message: string(data.Status),
// 		Data:    &data,
// 	}
// 	jsonData, _ := json.MarshalIndent(res, "", "   ")
// 	if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if data.Status != NoneTrade {
// 		if err := self.writeCSV(self.fileName, data, records); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 	}

// 	//withdraw
// 	if data.Status == Sell {
// 		if data.CashFlow > 20 {
// 			cashFlow := math.Round(data.CashFlow*100) / 100
// 			account, err := self.bitkubAPI.GetFiatAccount(1, 1)
// 			if err != nil || account.Error != 0 {
// 				msg := fmt.Sprintf("%v %d", err, account.Error)
// 				log.Println(msg)
// 				res := models.Response[any]{
// 					Title:   "Rebalance Bitkub btc_thb",
// 					Status:  false,
// 					Code:    500,
// 					Message: fmt.Sprintf("error bitkub api get fiat account : %s", msg),
// 					Data:    nil,
// 				}
// 				jsonData, _ := json.MarshalIndent(res, "", "   ")
// 				if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 					log.Println(err)
// 					return err
// 				}
// 				return err
// 			}

// 			resWithdraw, err := self.bitkubAPI.Withdraw(bitkub.WithdrawReq{
// 				ID:  account.Result[0].ID,
// 				Amt: cashFlow,
// 			})
// 			if err != nil || resWithdraw.Error != 0 {
// 				msg := fmt.Sprintf("%v %d", err, resWithdraw.Error)
// 				res := models.Response[any]{
// 					Title:   "Rebalance Bitkub btc_thb",
// 					Status:  false,
// 					Code:    500,
// 					Message: fmt.Sprintf("error bitkub api withdraw : %s", msg),
// 					Data:    nil,
// 				}
// 				jsonData, _ := json.MarshalIndent(res, "", "   ")
// 				if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 					log.Println(err)
// 					return err
// 				}
// 			}

// 			// success
// 			data := ResponseWithdraw{
// 				DateTime: time.Now().In(self.loc),
// 				Name:     account.Result[0].Name,
// 				Bank:     account.Result[0].Bank,
// 				CashFlow: cashFlow,
// 				Fee:      resWithdraw.Result.Fee,
// 				Receive:  resWithdraw.Result.Rec,
// 			}
// 			res := models.Response[ResponseWithdraw]{
// 				Title:   "Rebalance Bitkub btc_thb",
// 				Status:  true,
// 				Code:    200,
// 				Message: string(Withdraw),
// 				Data:    &data,
// 			}
// 			jsonData, _ := json.MarshalIndent(res, "", "   ")
// 			if err := self.lineAPI.SendMessage(string(jsonData)); err != nil {
// 				log.Println(err)
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
