package harrisonbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func getPERatioReport(stockID string, date string) (PERatioReport, error) {
	report := &ExchangeReport{}
	// today := time.Now().Format("20060102")
	// date := time.Now().AddDate(0, 0, -1).Format("20060102")
	url := TWSEAPIURL + "BWIBBU?response=json&date=" + date + "&stockNo=" + stockID

	res, err := http.Get(url)
	if err != nil {
		return PERatioReport{}, err
	}
	if res.StatusCode != http.StatusOK {
		return PERatioReport{}, errors.New("unexpected status" + res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(report); err != nil {
		fmt.Println("could not decode request body", err)
		return PERatioReport{}, err
	}

	if report.Status == "很抱歉，沒有符合條件的資料!" {
		return PERatioReport{}, errors.New("Wrong Stock ID")
	}

	if len(report.Data) == 0 {
		date := time.Now().AddDate(0, 0, -1).Format("20060102")
		return getPERatioReport(stockID, date)
	}

	index := len(report.Data) - 1
	existed, stockInfo := getStockInfo(stockID)
	var stockName string
	if existed {
		stockName = stockInfo.Name
	} else {
		stockName = ""
	}

	perReport := &PERatioReport{
		StockID:       stockID,
		StockName:     stockName,
		Date:          report.Data[index][0],
		DividendYield: report.Data[index][1],
		PER:           report.Data[index][3],
		PBR:           report.Data[index][4],
	}

	return *perReport, nil
}

func getStockPriceInfo(stockID string) (StockPriceInfo, error) {
	resp := &StockPriceResponse{}
	url := MISTWSEAPIURL + "stock/api/getStockInfo.jsp?ex_ch=tse_" + stockID + ".tw"
	res, err := http.Get(url)
	if err != nil {
		return StockPriceInfo{}, err
	}

	if res.StatusCode != http.StatusOK {
		return StockPriceInfo{}, errors.New("unexpected status" + res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		fmt.Println("could not decode request body", err)
		return StockPriceInfo{}, err
	}

	if len(resp.Message) == 0 {
		return StockPriceInfo{}, errors.New("Wrong Stock ID")
	}

	monthAvgPrice, err := getMonthAvgPrice(stockID, time.Now().Format("20060102"))
	if err != nil {
		monthAvgPrice = ""
	}

	info := &StockPriceInfo{
		StockID:               resp.Message[0].StockID,
		Name:                  resp.Message[0].Name,
		CurrentPrice:          resp.Message[0].CurrentPrice,
		OpeningPrice:          resp.Message[0].OpeningPrice,
		HighestPrice:          resp.Message[0].HighestPrice,
		LowestPrice:           resp.Message[0].LowestPrice,
		YesterdayClosingPrice: resp.Message[0].YesterdayClosingPrice,
		MonthAvgPrice:         monthAvgPrice,
	}

	return *info, nil
}

func getMonthAvgPrice(stockID string, date string) (string, error) {
	report := &ExchangeReport{}
	url := TWSEAPIURL + "STOCK_DAY_AVG?response=json&date=" + date + "&stockNo=" + stockID

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New("unexpected status" + res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(report); err != nil {
		fmt.Println("could not decode request body", err)
		return "", err
	}

	if report.Status == "很抱歉，沒有符合條件的資料!" {
		return "", errors.New("Wrong Stock ID")
	}

	if len(report.Data) == 0 {
		date := time.Now().AddDate(0, 0, -1).Format("20060102")
		return getMonthAvgPrice(stockID, date)
	}

	return report.Data[1][1].(string), nil
}

func getETFNetWorth(publisher string, stockID string) (ETFNetWorthInfo, error) {
	var etfInfo ETFNetWorthInfo
	etfNetWorth := &ETFNetWorth{}
	url := MISTWSEAPIURL + "stock/data/all_etf.txt"
	res, err := http.Get(url)
	if err != nil {
		return etfInfo, err
	}
	if res.StatusCode != http.StatusOK {
		return etfInfo, errors.New("unexpected status" + res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(etfNetWorth); err != nil {
		fmt.Println("could not decode request body", err)
		return etfInfo, err
	}

	var index int
	switch publisher {
	case "凱基":
		index = 0
	case "街口":
		index = 1
	case "台新":
		index = 2
	case "第一金":
		index = 3
	case "新光":
		index = 4
	case "中信":
		index = 5
	case "統一":
		index = 6
	case "FT":
		index = 7
	case "富邦":
		index = 8
	case "元大":
		index = 9
	case "國泰":
		index = 10
	case "永豐":
		index = 11
	case "FH":
		index = 12
	case "群益":
		index = 13
	case "兆豐":
		index = 14
	}

	for _, etf := range etfNetWorth.ETFPublisherList[index].ETFList {
		if etf.StockID == stockID {
			etfInfo = etf
			break
		}
	}

	return etfInfo, nil
}
