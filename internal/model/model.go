package model

type Message struct {
	ChatId      int64        `json:"chat_id"`
	Text        string       `json:"text"`
	From        From         `json:"from"`
	ParseMode   string       `json:"parse_mode,omitempty"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
}

type From struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type ReplyMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	IsPersistent    bool               `json:"is_persistent"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type Update struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

type Data struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}
