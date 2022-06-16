package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	tm "task-manager"
)

var mainKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Show tasks"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Add task"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Change tasks"),
	),
)

var showKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("All"),
		tg.NewKeyboardButton("Uncompleted"),
		tg.NewKeyboardButton("Overdue"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Return"),
	),
)

var changeKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Mark task"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Change task"),
		tg.NewKeyboardButton("Delete task"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Return"),
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
	var current tm.Task

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tg.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			msg.Text = "Hi! Select action type:"
			msg.ReplyMarkup = mainKeyboard
		case "Return":
			msg.Text = "Main menu:"
			msg.ReplyMarkup = mainKeyboard
		case "Show tasks":
			msg.Text = "Select view mode:"
			msg.ReplyMarkup = showKeyboard
		case "All":
			msg.Text = sl.ShowAll()
			msg.ReplyMarkup = showKeyboard
		case "Uncompleted":
			msg.Text = sl.ShowUncompleted()
			msg.ReplyMarkup = showKeyboard
		case "Overdue":
			msg.Text = sl.ShowOverdue()
			msg.ReplyMarkup = showKeyboard
		case "Change tasks":
			msg.Text = "Select task:"
			taskKeyboard := tg.ReplyKeyboardMarkup{ResizeKeyboard: true}
			for _, v := range *sl {
				taskKeyboard.Keyboard = append(taskKeyboard.Keyboard,
					tg.NewKeyboardButtonRow(tg.NewKeyboardButton(v.Title)))
			}
			msg.ReplyMarkup = taskKeyboard
		case "Mark task":
			sl.Mark(current.Title)
			msg.Text = current.Title + " completed!"
			msg.ReplyMarkup = mainKeyboard
		case "Delete task":
			sl.Delete(current.Title)
			msg.Text = current.Title + " deleted!"
			msg.ReplyMarkup = mainKeyboard
		default:
			// Find search task by title
			// It returns "undefined" if the task is not found
			current = sl.Find(update.Message.Text)

			if current.Title != "undefined" {
				// Print current task
				msg.Text = tm.Show(0, current)
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}

				msg.Text = "Select action:"
				msg.ReplyMarkup = changeKeyboard
			} else {
				msg.Text = current.Title
				msg.ReplyMarkup = mainKeyboard
			}
		}

		sl.Save("db.json")
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
