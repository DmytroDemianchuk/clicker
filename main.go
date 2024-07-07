package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var clickCount int // Змінна для зберігання кількості кліків

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	botToken := os.Getenv("TG_API_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TG_API_BOT_TOKEN is not set")
	}

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	// Обробник команди "start"
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)

		keyboard := tu.Keyboard(
			tu.KeyboardRow(
				tu.KeyboardButton("клік"),
				tu.KeyboardButton("Cancel"),
			),
		)

		message := tu.Message(
			chatID,
			"Привєт, шо тобі? Салатіку?",
		).WithReplyMarkup(keyboard)

		_, _ = bot.SendMessage(message)

	}, th.CommandEqual("start"))

	// Обробник для кнопки "клік"
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message == nil {
			return
		}

		chatID := tu.ID(update.Message.Chat.ID)

		clickCount++ // Збільшуємо кількість кліків

		_, _ = bot.SendMessage(tu.Message(chatID, fmt.Sprintf("Клік! Кількість кліків: %d", clickCount)))
	}, th.TextEqual("клік"))

	// Обробник для кнопки "Cancel"
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message == nil {
			return
		}

		chatID := tu.ID(update.Message.Chat.ID)

		// Відправлення повідомлення з видаленням клавіатури
		message := tu.Message(
			chatID,
			"Keyboard removed. Щоб вернути клавіатуру натисніть: /start",
		).WithReplyMarkup(tu.ReplyKeyboardRemove())

		_, _ = bot.SendMessage(message)
	}, th.TextEqual("Cancel"))

	bh.Start()
}
