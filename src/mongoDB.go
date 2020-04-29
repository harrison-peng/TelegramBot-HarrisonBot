package harrisonbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CheckUser is to check the user is existed in DB or not
func CheckUser(userID string, firstName string, lastName string) error {
	existed, _ := getUser(userID)
	if !existed {
		if _, err := insertUser(userID, firstName, lastName); err != nil {
			return err
		}
	}
	return nil
}

func connectDB() (mongo.Client, error) {
	mongoDBURL := "mongodb://" + MongoDBUser + ":" + MongoDBPassword + "@" + MongoDBDomain + MongoDBName
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoDBURL),
		options.Client().SetRetryWrites(false),
	)
	if err != nil {
		return mongo.Client{}, err
	}

	return *client, nil
}

func insertUser(userID string, firstName string, lastName string) (interface{}, error) {
	client, err := connectDB()
	if err != nil {
		return "", err
	}
	userCollection := client.Database(MongoDBName).Collection("User")
	value := bson.M{
		"_id":       userID,
		"FirstName": firstName,
		"LastName":  lastName,
		"StockList": []string{},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := userCollection.InsertOne(ctx, value)
	if err != nil {
		return "", err
	}

	return res.InsertedID, nil
}

func getUser(userID string) (bool, User) {
	user := &User{}

	client, err := connectDB()
	if err != nil {
		return false, *user
	}

	userCollection := client.Database(MongoDBName).Collection("User")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := bson.M{"_id": userID}
	if err := userCollection.FindOne(ctx, query).Decode(&user); err != nil {
		fmt.Println(err)
		return false, *user
	}
	// fmt.Printf("%+v\n", user)

	return true, *user
}

func getSession(userID string) string {
	_, user := getUser(userID)
	return user.Session
}

func updateSession(userID string, session string) error {
	client, err := connectDB()
	if err != nil {
		return err
	}

	userCollection := client.Database(MongoDBName).Collection("User")
	filter := bson.M{
		"_id": userID,
	}
	update := bson.M{
		"$set": bson.M{
			"session": session,
		},
	}

	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	}
	return errors.New("No matched user")
}

func getStockInfo(stockID string) (bool, Stock) {
	stockInfo := &Stock{}
	client, err := connectDB()
	if err != nil {
		return false, *stockInfo
	}

	stockCollection := client.Database(MongoDBName).Collection("StockInfo")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := bson.M{"_id": stockID}
	if err := stockCollection.FindOne(ctx, query).Decode(&stockInfo); err != nil {
		fmt.Println(err)
		return false, *stockInfo
	}
	// fmt.Printf("%+v\n", stockInfo)
	return true, *stockInfo
}

func addStock(userID string, stock Stock) error {
	client, err := connectDB()
	if err != nil {
		return err
	}

	userCollection := client.Database(MongoDBName).Collection("User")
	filter := bson.M{
		"_id": userID,
	}
	update := bson.M{
		"$push": bson.M{
			"StockList": bson.M{
				"_id":       stock.StockID,
				"name":      stock.Name,
				"type":      stock.Type,
				"stockType": stock.StockType,
			},
		},
	}

	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	}
	return errors.New("No matched user")
}

func removeStock(userID string, stockID string) error {
	client, err := connectDB()
	if err != nil {
		return err
	}

	userCollection := client.Database(MongoDBName).Collection("User")
	filter := bson.M{
		"_id": userID,
	}
	update := bson.M{
		"$pull": bson.M{
			"StockList": bson.M{
				"_id": stockID,
			},
		},
	}

	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	}
	return errors.New("No matched user")
}
