package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	m "tgbot_internship_vk/internal/model"

	"time"
)

type App struct {
	token  string
	runing bool
}

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

func (a *App) sendMessage(messageBytes []byte) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", a.token)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *App) getUpdates(offset int64) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", a.token, offset)
}

func (a *App) Run() {
	var lastId int64 = 0

	log.Println("запуск приложения")

	for a.runing {
		time.Sleep(time.Duration(1) * time.Second)

		url := a.getUpdates(lastId)

		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Printf("API Telegram вернул ошибку с кодом: %d\n", response.StatusCode)
			continue
		}

		data := m.Data{}

		if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
			fmt.Println(err)
			continue
		}
		for _, update := range data.Result {

			if lastId < update.UpdateID {
				log.Printf("Пользователь: %s, текст: %s", update.Message.From.Username, update.Message.Text)
				message := m.Message{}

				if update.Message.Text == "/start" {
					message = m.Message{
						ChatId:      update.Message.From.Id,
						Text:        update.Message.From.Username + ", добро пожаловать!",
						ReplyMarkup: &keyboard,
					}
				} else if strings.HasPrefix(update.Message.Text, "Кнопка") {
					message = m.Message{
						ChatId:      update.Message.From.Id,
						Text:        "Вы нажали: " + update.Message.Text,
						ReplyMarkup: &keyboard,
					}
				} else {
					message = m.Message{
						ChatId:      update.Message.From.Id,
						Text:        "Вы ввели: " + update.Message.Text,
						ReplyMarkup: &keyboard,
					}
				}

				messageBytes, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
					continue
				}
				lastId = update.UpdateID
				err = a.sendMessage(messageBytes)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

func (a *App) Shutdown() {
	a.runing = false
}
