package harrisonbot

// User is the user object of the mongo DB
type User struct {
	UserID    string   `bson:"_id"`
	FirstName string   `bson:"FirstName"`
	LastName  string   `bson:"LastName"`
	StockList []string `bson:"StockList"`
	Session   string   `bson:"session"`
}

// Stock is the stock info object of the mongo DB
type Stock struct {
	STockID   string `bson:"_id"`
	Name      string `bson:"name"`
	Type      string `bson:"type"`
	StockType string `bson:"stockType"`
	Publisher string `bson:"publisher"`
}
