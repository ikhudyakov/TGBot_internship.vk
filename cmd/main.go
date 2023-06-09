package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	h "tgbot_internship_vk/internal/handler"

	"github.com/joho/godotenv"
)

var app App

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("файл .env не найден")
	}
}

func main() {
	token, exists := os.LookupEnv("TOKEN")

	if !exists {
		log.Fatal("в файле .env не указаны переменные окружения: TOKEN")
	}

	handler := h.Handler{}

	app = App{
		token:   token,
		handler: handler,
	}

	go func() {
		app.runing = true
		app.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("приложение останавливается")

	app.Shutdown()
}
