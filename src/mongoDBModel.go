package harrisonbot

// User is the user object of the mongo DB
type User struct {
	UserID    string   `bson:"_id"`
	FirstName string   `bson:"FirstName"`
	LastName  string   `bson:"LastName"`
	StockList []string `bson:"StockList"`
	Session   string   `bson:"session"`
}
