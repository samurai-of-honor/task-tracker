package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	tm "task-manager"
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Show"),
		tgbotapi.NewKeyboardButton("Change"),
	),
)

var showKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("All"),
		tgbotapi.NewKeyboardButton("Uncompleted"),
		tgbotapi.NewKeyboardButton("Overdue"),
	),
)

func Auth() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

func main() {
	bot := Auth()

	sl := tm.Create()
	sl.Load("db.json")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			msg.Text = "Hi! Select action:"
			msg.ReplyMarkup = mainKeyboard
		case "Show":
			msg.Text = "Select view mode:"
			msg.ReplyMarkup = showKeyboard
		case "All":
			msg.Text = sl.ShowAll()
			msg.ReplyMarkup = mainKeyboard
		case "Uncompleted":
			msg.Text = sl.ShowUncompleted()
			msg.ReplyMarkup = mainKeyboard
		case "Overdue":
			msg.Text = sl.ShowOverdue()
			msg.ReplyMarkup = mainKeyboard
		default:
			msg.Text = "undefined"
			msg.ReplyMarkup = mainKeyboard
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
