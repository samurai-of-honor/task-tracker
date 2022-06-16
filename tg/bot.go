package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	tm "task-manager"
)

var bot *tg.BotAPI

func Auth() {
	var err error
	bot, err = tg.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func Send(msg tg.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func GetText(updates tg.UpdatesChannel) string {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		return update.Message.Text
	}
	return "Unknown command😕"
}

func main() {
	Auth()

	sl := tm.Create()
	sl.Load("db.json")
	var curr tm.Task

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
			msg.Text = "Hi👋 Select action type:"
			msg.ReplyMarkup = mainKeyboard
		case "🔙 Return":
			msg.Text = "🗄 Main menu:"
			msg.ReplyMarkup = mainKeyboard
		case "🖼 Show tasks":
			msg.Text = "🔍 Select view mode:"
			msg.ReplyMarkup = showKeyboard
		case "🗃 All":
			msg.Text = sl.ShowAll()
			msg.ReplyMarkup = showKeyboard
		case "🛑 Uncompleted":
			msg.Text = sl.ShowUncompleted()
			msg.ReplyMarkup = showKeyboard
		case "‼️Overdue":
			msg.Text = sl.ShowOverdue()
			msg.ReplyMarkup = showKeyboard
		case "📌 Add task":
			mainKeyboard.OneTimeKeyboard = true // Hide main keyboard for writing
			msg.Text = "✏️ Enter title:"
			Send(msg)
			title := GetText(updates)
			msg.Text = "📝 Enter description:"
			Send(msg)
			desc := GetText(updates)
			msg.Text = "⏰ Enter deadline in this format \"01-01-2022 12:00\":"
			Send(msg)
			deadline := GetText(updates)

			mainKeyboard.OneTimeKeyboard = false // Show main keyboard
			msg.ReplyMarkup = mainKeyboard
			if err := sl.Add(title, desc, deadline); err != nil {
				msg.Text = "Deadline error ❌"
			} else {
				msg.Text = title + " added ✅"
				Send(msg)
				continue
			}
		case "🛠 Change tasks":
			msg.Text = "🗂 Select task:"

			taskKeyboard := tg.ReplyKeyboardMarkup{ResizeKeyboard: true}
			for i, v := range *sl {
				if i%2 == 0 {
					taskKeyboard.Keyboard = append(taskKeyboard.Keyboard,
						tg.NewKeyboardButtonRow(tg.NewKeyboardButton("🔹 "+v.Title)))
				} else {
					taskKeyboard.Keyboard[i/2] = append(taskKeyboard.Keyboard[i/2],
						tg.NewKeyboardButton("🔸 "+v.Title))
				}
			}
			taskKeyboard.Keyboard = append(taskKeyboard.Keyboard,
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton("🔙 Return")))

			msg.ReplyMarkup = taskKeyboard
		case "✅ Mark task":
			sl.Mark(curr.Title)
			msg.Text = curr.Title + " completed ✅"
			msg.ReplyMarkup = mainKeyboard
		case "🔧 Change task":
			msg.Text = "⚙️ Which option do you want to change?"
			optionKeyboard.OneTimeKeyboard = true
			msg.ReplyMarkup = optionKeyboard
		case "✏️ Title":
			msg.Text = "✏️ Enter new title:"
			Send(msg)
			title := GetText(updates)
			sl.Change(curr.Title, title, curr.Description, curr.Deadline)
			msg.Text = "Title changed ✅"
			msg.ReplyMarkup = changeKeyboard
		case "📝 Description":
			msg.Text = "📝 Enter new description:"
			Send(msg)
			desc := GetText(updates)
			sl.Change(curr.Title, curr.Title, desc, curr.Deadline)
			msg.Text = "Description changed ✅"
			msg.ReplyMarkup = changeKeyboard
		case "⏰ Deadline":
			msg.Text = "⏰ Enter new deadline:"
			Send(msg)
			deadline := GetText(updates)
			sl.Change(curr.Title, curr.Title, curr.Description, deadline)
			msg.Text = "Deadline changed ✅"
			msg.ReplyMarkup = changeKeyboard
		case "🗑 Delete task":
			sl.Delete(curr.Title)
			msg.Text = curr.Title + " deleted ✅"
			msg.ReplyMarkup = mainKeyboard
		default:
			// Find search task by title
			// It returns "undefined" if the task is not found
			curr = sl.Find(update.Message.Text)

			if curr.Title != "undefined" {
				// Print curr task
				msg.Text = tm.Show(0, curr)
				Send(msg)

				msg.Text = "🎬 Select action:"
				msg.ReplyMarkup = changeKeyboard
			} else {
				msg.Text = curr.Title
				msg.ReplyMarkup = mainKeyboard
			}
		}

		sl.Save("db.json")
		Send(msg)
	}
}
