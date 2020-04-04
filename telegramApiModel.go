package main

// Update is the result of Update
// https://core.telegram.org/bots/api#update
type Update struct {
	UpdateID      int           `json:"update_id"`
	Message       Message       `'json:"message"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

// Message is the message of Result
type Message struct {
	MessageID int     `json:"message_id"`
	Date      int     `json:"date"`
	From      From    `json:"from"`
	Chat      Chat    `json:"chat"`
	Text      string  `json:"text"`
	Sticker   Sticker `json:"sticker"`
}

// From is the sender of the message
type From struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	LanguageCode string `json:"language_code"`
}

// Chat is the chat of the message
type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
}

// Sticker is the object od sticker
type Sticker struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Emoji        string `json:"emoji"`
	SetName      string `json:"set_name"`
	IsAnimated   bool   `json:"is_animated"`
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Thumb        Thumb  `json:"thumb"`
}

// Thumb is the thumb of the sticker
type Thumb struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// CallbackQuery is the callback query object of the Update
type CallbackQuery struct {
	ID            string  `json:"id"`
	From          From    `json:"from"`
	Message       Message `json:"message"`
	ChatInstance  string  `json:"chat_instance"`
	Data          string  `json:"data"`
	GameShortName string  `json:"game_short_name"`
}

// SendMessageRequest is the request body of sending message
// https://core.telegram.org/bots/api#sendmessage
type SendMessageRequest struct {
	ChatID           int          `json:"chat_id"`
	Text             string       `json:"text"`
	ParseMode        string       `json:"parse_mode,omitempty"`
	ReplyToMessageID int          `json:"reply_to_message_id,omitempty"`
	ReplyMarkup      *ReplyMarkup `json:"reply_markup,omitempty"`
}

// ReplyMarkup is the ReplyMarkup of the SendMessageRequest
type ReplyMarkup struct {
	Keyboard        *[][]ReplyKeyboardButton  `json:"keyboard,omitempty"`
	ResizeKeyboard  *bool                     `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard *bool                     `json:"one_time_keyboard,omitempty"`
	Selective       *bool                     `json:"selective,omitempty"`
	InlineKeyboard  *[][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
	RemoveKeyboard  *bool                     `json:"remove_keyboard,omitempty"`
	ForceReply      *bool                     `json:"force_reply,omitempty"`
}

// ReplyKeyboardButton is the reply button object of the ReplyMarkup
type ReplyKeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
	RequestPoll     *struct {
		Type string `json:"type"`
	} `json:"request_poll,omitempty"`
}

// InlineKeyboardButton is the inline button object of the ReplyMarkup
type InlineKeyboardButton struct {
	Text                         string  `json:"text"`
	URL                          *string `json:"url,omitempty"`
	CallbackData                 *string `json:"callback_data,omitempty"`
	SwitchInlineQuery            *string `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat *string `json:"switch_inline_query_current_chat,omitempty"`
	Pay                          *bool   `json:"pay,omitempty"`
}

// SendStickerRequest is the request body of sending sticker
type SendStickerRequest struct {
	ChatID  int    `json:"chat_id"`
	Sticker string `json:"sticker"`
}

// EditMessageReplyMarkup is the object of the edit reply markup message
type EditMessageReplyMarkup struct {
	ChatID      int         `json:"chat_id"`
	MessageID   int         `json:"message_id"`
	Text        string      `json:"text"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}
