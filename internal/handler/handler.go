package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	m "tgbot_internship_vk/internal/model"
)

type Handler struct{}

var keyboard = m.ReplyMarkup{
	Keyboard: [][]m.KeyboardButton{
		{
			{
				Text: "Кнопка 1",
			},
			{
				Text: "Кнопка 2",
			},
			{
				Text: "Кнопка 3",
			},
			{
				Text: "Кнопка 4",
			},
		},
		{
			{
				Text: "Кнопка 5",
			},
			{
				Text: "Кнопка 6",
			},
		},
	},
	ResizeKeyboard:  true,
	OneTimeKeyboard: true,
	IsPersistent:    true,
}

func (h *Handler) Control(response *http.Response, lastId *int64) (message *m.Message, err error) {

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("API Telegram вернул ошибку с кодом: %d", response.StatusCode)
		return nil, err
	}

	data := m.Data{}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	for _, update := range data.Result {

		if *lastId < update.UpdateID {
			log.Printf("Пользователь: %s, текст: %s", update.Message.From.Username, update.Message.Text)

			if update.Message.Text == "/start" {
				message = &m.Message{
					ChatId:      update.Message.From.Id,
					Text:        update.Message.From.Username + ", добро пожаловать!",
					ReplyMarkup: &keyboard,
				}
			} else if strings.HasPrefix(update.Message.Text, "Кнопка") {
				message = &m.Message{
					ChatId:      update.Message.From.Id,
					Text:        "Вы нажали: " + update.Message.Text,
					ReplyMarkup: &keyboard,
				}
			} else {
				message = &m.Message{
					ChatId:      update.Message.From.Id,
					Text:        "Вы ввели: " + update.Message.Text,
					ReplyMarkup: &keyboard,
				}
			}
			*lastId = update.UpdateID
		}
	}
	return message, nil
}
