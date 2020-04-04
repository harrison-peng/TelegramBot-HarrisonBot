package main

// ExchangeReport is the API response of exchange report
// https://www.twse.com.tw/exchangeReport/BWIBBU?response=json&date={date}&stockNo={stockNum}
type ExchangeReport struct {
	Status string          `json:"stat"`
	Title  string          `json:"title"`
	Fields []string        `json:"fields"`
	Data   [][]interface{} `json:"data"`
}

// PERatioReport is the PE Ratio Report of the Exchange Report
type PERatioReport struct {
	StockID       string
	StockName     string
	Date          interface{}
	DividendYield interface{}
	PER           interface{}
	PBR           interface{}
}

// StockPriceResponse is the API response of stock info
// https://mis.twse.com.tw/stock/index.jsp
// https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_{stockNum}.tw
type StockPriceResponse struct {
	QueryTime struct {
		Date string `json:"sysDate"`
		Time string `json:"sysTime"`
	} `json:"queryTime"`
	Status  string `json:"rtmessage"`
	Message []struct {
		StockID               string `json:"c"`
		Name                  string `json:"n"`
		CurrentPrice          string `json:"z"`
		CurrentQuantity       string `json:"tv"`
		CumulatedQuantity     string `json:"v"`
		OpeningPrice          string `json:"o"`
		HighestPrice          string `json:"h"`
		LowestPrice           string `json:"l"`
		YesterdayClosingPrice string `json:"y"`
	} `json:"msgArray"`
}

// StockPriceInfo is the object of the price info
type StockPriceInfo struct {
	StockID               string
	Name                  string
	CurrentPrice          string
	OpeningPrice          string
	HighestPrice          string
	LowestPrice           string
	YesterdayClosingPrice string
	MonthAvgPrice         string
}

// StockInfo is the infomation of the Stock in Mongo DB
type StockInfo struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Type string `bson:"type"`
}
