package harrisonbot

import "os"

const (
	// APIURL is the api url of the telegram bot
	APIURL = "https://api.telegram.org/"
	// TWSEAPIURL is the api url of TWSE
	TWSEAPIURL = "https://www.twse.com.tw/exchangeReport/"
	// MISTWSEAPIURL is the api url of mis TWSE
	MISTWSEAPIURL = "https://mis.twse.com.tw/stock/api/getStockInfo.jsp"
)

// Token is the token of the telegram bot
var Token = os.Getenv("TELEGRAM_API_TOKEN")

// MongoDBUser is the user name of the mongo DB
var MongoDBUser = os.Getenv("MONGO_DB_USER")

// MongoDBPassword is the password of the mongo DB
var MongoDBPassword = os.Getenv("MONGO_DB_PWD")

// MongoDBName is the DB name of the mongo DB
var MongoDBName = os.Getenv("MONGO_DB_NAME")

// MongoDBDomain is the domain name of the mongo DB
var MongoDBDomain = os.Getenv("MONGO_DB_DOMAIN")
