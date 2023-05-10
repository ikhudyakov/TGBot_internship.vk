package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	h "tgbot_internship_vk/internal/handler"

	"time"
)

type App struct {
	token   string
	runing  bool
	handler h.Handler
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

		message, err := a.handler.Control(response, &lastId)
		if err != nil {
			log.Println(err)
			continue
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
			continue
		}

		err = a.sendMessage(messageBytes)
		if err != nil {
			log.Println(err)
			continue
		}
	}

}

func (a *App) Shutdown() {
	a.runing = false
}
