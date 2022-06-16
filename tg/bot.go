package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	tm "task-manager"
)

var mainKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Show"),
		tg.NewKeyboardButton("Change"),
	),
)

var showKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("All"),
		tg.NewKeyboardButton("Uncompleted"),
		tg.NewKeyboardButton("Overdue"),
	),
)

var changeKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Mark task"),
		tg.NewKeyboardButton("Add task"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Change task"),
		tg.NewKeyboardButton("Delete task"),
	),
)

func Auth() *tg.BotAPI {
	bot, err := tg.NewBotAPI(token)
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

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var curr string
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tg.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			msg.Text = "Hi! Select action type:"
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
		case "Change":
			msg.Text = "Select task:"
			taskKeyboard := tg.ReplyKeyboardMarkup{}
			for _, v := range *sl {
				taskKeyboard.Keyboard = append(taskKeyboard.Keyboard,
					tg.NewKeyboardButtonRow(tg.NewKeyboardButton(v.Title)))
			}
			msg.ReplyMarkup = taskKeyboard
		case "Mark task":
			sl.Mark(curr)
			msg.Text = curr + " completed!"
			msg.ReplyMarkup = mainKeyboard
		default:
			curr = sl.Find(update.Message.Text)
			if curr != "undefined" {
				msg.Text = "Select action:"
				msg.ReplyMarkup = changeKeyboard
			} else {
				msg.Text = curr
				msg.ReplyMarkup = mainKeyboard
			}
		}

		sl.Save("db.json")
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
