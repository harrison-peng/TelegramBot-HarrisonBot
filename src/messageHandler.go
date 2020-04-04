package harrisonbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HandleCallbackData is
func HandleCallbackData(callbackQuery CallbackQuery) error {
	if callbackQuery.Data == "/addStock" {
		if err := updateSession(strconv.Itoa(callbackQuery.Message.Chat.ID), "addStock"); err != nil {
			return err
		}
		if err := sendMessage(callbackQuery.Message.Chat.ID, "Please input the stock number:", "", ReplyMarkup{}); err != nil {
			return err
		}
	} else if callbackQuery.Data == "/removeStockPage" {
		_, user := getUser(strconv.Itoa(callbackQuery.Message.Chat.ID))
		if len(user.StockList) > 0 {
			var button [][]InlineKeyboardButton
			for _, stockID := range user.StockList {
				_, stockInfo := getStockInfo(stockID)
				callbackData := "/removeStock:" + stockInfo.ID
				button = append(button, []InlineKeyboardButton{
					InlineKeyboardButton{
						Text:         stockInfo.ID + " " + stockInfo.Name,
						CallbackData: &callbackData,
					},
				})
			}
			replyMarkup := ReplyMarkup{
				InlineKeyboard: &button,
			}
			if err := editMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Choose the stock to remove:", replyMarkup); err != nil {
				return err
			}
		} else {
			content := "*Reminder:*\nYou don't have any stock in your list."
			if err := sendMessage(callbackQuery.Message.Chat.ID, content, "Markdown", ReplyMarkup{}); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(callbackQuery.Data, "/removeStock:") {
		stockID := strings.Replace(strings.Split(callbackQuery.Data, "/removeStock:")[1], " ", "", -1)
		if err := removeStock(strconv.Itoa(callbackQuery.Message.Chat.ID), stockID); err != nil {
			return err
		}
		content := "Remove " + stockID + " successed!"
		if err := sendMessage(callbackQuery.Message.Chat.ID, content, "", ReplyMarkup{}); err != nil {
			return err
		}
	} else if callbackQuery.Data == "/getMyStocksPrice" {
		existed, user := getUser(strconv.Itoa(callbackQuery.Message.Chat.ID))
		if existed {
			stockList := user.StockList
			for _, stockID := range stockList {
				info, _ := getStockPriceInfo(stockID)
				if err := sendStockInfoMessage(callbackQuery.Message.Chat.ID, info); err != nil {
					return err
				}
			}
		}
	} else if callbackQuery.Data == "/getMyStocksPER" {
		existed, user := getUser(strconv.Itoa(callbackQuery.Message.Chat.ID))
		if existed {
			stockList := user.StockList
			for _, stockID := range stockList {
				report, err := getPERatioReport(stockID, time.Now().Format("20060102"))
				if err != nil {
					return err
				}

				if err := sendStockPERMessage(callbackQuery.Message.Chat.ID, stockID, report); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// HandleMessage is
func HandleMessage(message Message) error {
	if strings.HasPrefix(message.Text, "/") {
		if err := updateSession(strconv.Itoa(message.Chat.ID), ""); err != nil {
			return err
		}
	}

	session := getSession(strconv.Itoa(message.Chat.ID))
	if session != "" {
		if session == "addStock" {
			existed, stockInfo := getStockInfo(message.Text)
			if existed {
				if err := addStock(strconv.Itoa(message.Chat.ID), stockInfo.ID); err != nil {
					return err
				}
				if err := updateSession(strconv.Itoa(message.Chat.ID), ""); err != nil {
					return err
				}
				content := "Insert " + stockInfo.ID + " " + stockInfo.Name + " successed!"
				if err := sendMessage(message.Chat.ID, content, "", ReplyMarkup{}); err != nil {
					return err
				}
			} else {
				if err := sendMessage(message.Chat.ID, "Wrong stock ID, please type again:", "", ReplyMarkup{}); err != nil {
					return err
				}
			}
		} else if session == "searchStock" {
			existed, stockInfo := getStockInfo(message.Text)
			if existed {
				info, err := getStockPriceInfo(stockInfo.ID)
				if err != nil {
					return err
				}

				if err := sendStockInfoMessage(message.Chat.ID, info); err != nil {
					return err
				}
			} else {
				if err := sendMessage(message.Chat.ID, "Wrong stock ID, please type again:", "", ReplyMarkup{}); err != nil {
					return err
				}
			}
		}
	} else {
		if message.Text != "" {
			fmt.Println("Text: ", message.Text)
			if message.Text == "/start" {
				if err := sendReplyKeyboardMessage(message.Chat.ID); err != nil {
					return err
				}
				return nil
			} else if message.Text == "/settings" {
				if err := sendSettingKeyboardMessage(message.Chat.ID); err != nil {
					return err
				}
			} else if message.Text == "MyStocks" {
				userID := strconv.Itoa(message.Chat.ID)
				existed, user := getUser(userID)
				if existed {
					if len(user.StockList) > 0 {
						if err := sendMyStockKeyboardMessage(message.Chat.ID); err != nil {
							return err
						}
					} else {
						content := "*Reminder:*\nYou don't have any stock in your list, you can add stocks to your list in /settings"
						if err := sendMessage(message.Chat.ID, content, "Markdown", ReplyMarkup{}); err != nil {
							return err
						}
					}
				}
			} else if message.Text == "SearchStock" {
				if err := updateSession(strconv.Itoa(message.Chat.ID), "searchStock"); err != nil {
					return err
				}

				if err := sendMessage(message.Chat.ID, "Please input the stock number:", "", ReplyMarkup{}); err != nil {
					return err
				}
			} else {
				if err := sendMessage(message.Chat.ID, "此功能尚未開發！", "", ReplyMarkup{}); err != nil {
					return err
				}
			}

		} else if message.Sticker.FileID != "" {
			if err := sendSticker(message.Chat.ID, message.Sticker.FileID); err != nil {
				return err
			}
		}
	}

	return nil
}

func sendReplyKeyboardMessage(chatID int) error {
	trueVal := true
	replyMarkup := &ReplyMarkup{
		Keyboard: &[][]ReplyKeyboardButton{
			[]ReplyKeyboardButton{
				ReplyKeyboardButton{
					Text: "MyStocks",
				},
				ReplyKeyboardButton{
					Text: "SearchStock",
				},
			},
		},
		ResizeKeyboard:  &trueVal,
		OneTimeKeyboard: &trueVal,
	}

	if err := sendMessage(chatID, "Hi, please use the buttons below to get the information of the stocks.", "", *replyMarkup); err != nil {
		return err
	}
	return nil
}

func sendSettingKeyboardMessage(chatID int) error {
	addStockCallbackData := "/addStock"
	removeStockCallbackData := "/removeStockPage"
	replyMarkup := &ReplyMarkup{
		InlineKeyboard: &[][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				InlineKeyboardButton{
					Text:         "Add Stock To MyStocks",
					CallbackData: &addStockCallbackData,
				},
			},
			[]InlineKeyboardButton{
				InlineKeyboardButton{
					Text:         "Remove Stock From MyStocks",
					CallbackData: &removeStockCallbackData,
				},
			},
		},
	}
	// fmt.Printf("%+v\n", *replyMarkup)
	if err := sendMessage(chatID, "Settings:", "", *replyMarkup); err != nil {
		return err
	}
	return nil
}

func sendMyStockKeyboardMessage(chatID int) error {
	priceCallbackData := "/getMyStocksPrice"
	perCallbackData := "/getMyStocksPER"
	replyMarkup := &ReplyMarkup{
		InlineKeyboard: &[][]InlineKeyboardButton{
			[]InlineKeyboardButton{
				InlineKeyboardButton{
					Text:         "Price Information",
					CallbackData: &priceCallbackData,
				},
			},
			[]InlineKeyboardButton{
				InlineKeyboardButton{
					Text:         "PE Ratio & Others",
					CallbackData: &perCallbackData,
				},
			},
		},
	}
	// fmt.Printf("%+v\n", *replyMarkup)
	if err := sendMessage(chatID, "Choose the information:", "", *replyMarkup); err != nil {
		return err
	}
	return nil
}

func sendStockInfoMessage(chatID int, info StockPriceInfo) error {
	var content bytes.Buffer
	content.WriteString("*" + info.StockID + " " + info.Name + "*\n")
	content.WriteString("*現價:* " + info.CurrentPrice + " 元\n")
	content.WriteString("*作收:* " + info.YesterdayClosingPrice + " 元\n")
	content.WriteString("*開盤:* " + info.OpeningPrice + " 元\n")
	content.WriteString("*最高:* " + info.HighestPrice + " 元\n")
	content.WriteString("*最低:* " + info.LowestPrice + " 元\n")
	content.WriteString("*月均:* " + info.MonthAvgPrice + " 元\n")

	if err := sendMessage(chatID, content.String(), "Markdown", ReplyMarkup{}); err != nil {
		return err
	}

	return nil
}

func sendStockPERMessage(chatID int, stockID string, report PERatioReport) error {
	var content bytes.Buffer
	content.WriteString("*" + stockID + " " + report.StockName + "*\n")
	content.WriteString("*日期:* " + report.Date.(string) + "\n")
	content.WriteString("*殖利率(%):* " + report.DividendYield.(string) + "\n")
	content.WriteString("*本益比:* " + report.PER.(string) + "\n")
	content.WriteString("*股價淨值比:* " + report.PBR.(string) + "\n")

	if err := sendMessage(chatID, content.String(), "Markdown", ReplyMarkup{}); err != nil {
		return err
	}

	return nil
}

// replyMessage takes a chatID and sends content to them
func sendMessage(chatID int, content string, parseMode string, replyMarkup ReplyMarkup) error {
	var reqBody SendMessageRequest
	if &replyMarkup == nil {
		reqBody = SendMessageRequest{
			ChatID:    chatID,
			Text:      content,
			ParseMode: parseMode,
		}
	} else {
		reqBody = SendMessageRequest{
			ChatID:      chatID,
			Text:        content,
			ParseMode:   parseMode,
			ReplyMarkup: &replyMarkup,
		}
	}

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	sendMessageURL := APIURL + Token + "/sendMessage"
	res, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func editMessage(chatID int, messageID int, content string, replyMarkup ReplyMarkup) error {
	reqBody := EditMessageReplyMarkup{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        content,
		ReplyMarkup: replyMarkup,
	}
	// fmt.Printf("%+v\n", reqBody)

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	sendMessageURL := APIURL + Token + "/editMessageText"
	res, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func sendSticker(chatID int, sticker string) error {
	reqBody := &SendStickerRequest{
		ChatID:  chatID,
		Sticker: sticker,
	}

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	sendMessageURL := APIURL + Token + "/sendSticker"
	res, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
