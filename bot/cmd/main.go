package main

import (
	"log"
	"todo/bot/pkg/adapter"
	"todo/bot/pkg/telegram"
	"todo/internal/db"
	"todo/internal/repo"
	"todo/internal/usecases"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tg.NewBotAPI("7999162539:AAGW9pvXKUUxAAEj-ykE5YCaFgRpslbgJik")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	database, _ := db.Connection()

	newRepo := repo.NewTaskRepositoryImpl(database)
	newUsecases := usecases.NewTaskServiceImpl(newRepo)
	newAdapter := adapter.NewBotAdapter(newUsecases)

	telegramBot := telegram.NewBot(bot, *newAdapter)

	if err := telegramBot.Start(); err != nil {
		log.Panic(err)
	}

}
