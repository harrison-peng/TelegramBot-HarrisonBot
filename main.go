package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	harrisonBot "github.com/hichyen1207/TelegramBot-HarrisonBot/src"
)

// Handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	fmt.Println(harrisonBot.APIURL)
	// Decode the JSON response body
	body := &harrisonBot.Update{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}
	fmt.Printf("%+v\n", body)

	if body.Message.MessageID == 0 {
		// Handle callback data
		if err := harrisonBot.HandleCallbackData(body.CallbackQuery); err != nil {
			fmt.Println("Error in handling callback data: ", err)
			return
		}
	} else {
		// Handle text message
		if err := harrisonBot.HandleMessage(body.Message); err != nil {
			fmt.Println("Error in handling message: ", err)
			return
		}
	}
}

// the main funtion starts our server on port 3000
func main() {
	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}
