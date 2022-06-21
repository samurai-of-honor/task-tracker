package main

import (
	"flag"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	tm "task-manager"
	. "task-manager/tg/localization"
)

var token = flag.String("t", "", "Token for access bot")
var bot *tg.BotAPI

func Auth() {
	var err error
	flag.Parse()
	bot, err = tg.NewBotAPI(*token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func TaskFile(id int64) string {
	sID := "./tg/taskBases/" + strconv.Itoa(int(id)) + ".json"
	return sID
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
	return UnknownCommand
}

func UpdateHandler() {
	var curr tm.Task

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		sl := tm.Create()
		sl.Load(TaskFile(update.Message.Chat.ID))

		msg := tg.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			msg.Text = StartText
			msg.ReplyMarkup = mainKeyboard
		case Return:
			msg.Text = MainMenu
			msg.ReplyMarkup = mainKeyboard
		case ShowTasks:
			msg.Text = ViewMode
			msg.ReplyMarkup = showKeyboard
		case AllTasks:
			text := sl.ShowAll()
			if text == "" {
				msg.Text = Nothing
			} else {
				msg.Text = text
			}
			msg.ReplyMarkup = showKeyboard
		case UncompletedTasks:
			text := sl.ShowUncompleted()
			if text == "" {
				msg.Text = Nothing
			} else {
				msg.Text = text
			}
			msg.ReplyMarkup = showKeyboard
		case OverdueTasks:
			text := sl.ShowOverdue()
			if text == "" {
				msg.Text = Nothing
			} else {
				msg.Text = text
			}
			msg.ReplyMarkup = showKeyboard
		case AddTask:
			msg.Text = SelectTask
			newTaskKeyboard := tg.ReplyKeyboardMarkup{OneTimeKeyboard: true}
			newTaskKeyboard.Keyboard = append(newTaskKeyboard.Keyboard,
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(NewTask)))
			msg.ReplyMarkup = newTaskKeyboard
		case NewTask:
			msg.Text = InTitle
			Send(msg)
			title := GetText(updates)
			msg.Text = InDesc
			Send(msg)
			desc := GetText(updates)
			msg.Text = InDLine
			Send(msg)
			deadline := GetText(updates)

			err := sl.Add(title, desc, deadline)
			if err == "deadline error" {
				msg.Text = ErrDLine
			} else if err == "title error" {
				msg.Text = ErrTitle
			} else {
				msg.Text = title + TaskAdded
			}
			msg.ReplyMarkup = mainKeyboard
		case ChangeTasks:
			msg.Text = SelectTask

			tasksKeyboard := tg.ReplyKeyboardMarkup{ResizeKeyboard: true}
			for i, v := range *sl {
				if i%2 == 0 {
					tasksKeyboard.Keyboard = append(tasksKeyboard.Keyboard,
						tg.NewKeyboardButtonRow(tg.NewKeyboardButton(BlueRhombus+v.Title)))
				} else {
					tasksKeyboard.Keyboard[i/2] = append(tasksKeyboard.Keyboard[i/2],
						tg.NewKeyboardButton(YellowRhombus+v.Title))
				}
			}
			tasksKeyboard.Keyboard = append(tasksKeyboard.Keyboard,
				tg.NewKeyboardButtonRow(tg.NewKeyboardButton(Return)))

			msg.ReplyMarkup = tasksKeyboard
		case MarkTask:
			sl.Mark(curr.Title)
			msg.Text = curr.Title + TaskCompleted
			msg.ReplyMarkup = mainKeyboard
		case ChangeTask:
			msg.Text = TaskOptions
			optionKeyboard.OneTimeKeyboard = true
			msg.ReplyMarkup = optionKeyboard
		case ChangeTitle:
			msg.Text = NewTitle
			Send(msg)
			title := GetText(updates)
			sl.Change(curr.Title, title, curr.Description, curr.Deadline)
			curr = sl.Find(title)

			msg.Text = ChangeTitle + TaskChanged
			msg.ReplyMarkup = changeKeyboard
		case ChangeDesc:
			msg.Text = NewDesc
			Send(msg)
			desc := GetText(updates)
			sl.Change(curr.Title, curr.Title, desc, curr.Deadline)
			msg.Text = ChangeDesc + TaskChanged
			msg.ReplyMarkup = changeKeyboard
		case ChangeDLIne:
			msg.Text = NewDLine
			Send(msg)
			deadline := GetText(updates)
			sl.Change(curr.Title, curr.Title, curr.Description, deadline)
			msg.Text = ChangeDLIne + TaskChanged
			msg.ReplyMarkup = changeKeyboard
		case DeleteTask:
			sl.Delete(curr.Title)
			msg.Text = curr.Title + TaskDeleted
			msg.ReplyMarkup = mainKeyboard
		default:
			// Find search task by title
			// It returns "undefined" if the task is not found
			curr = sl.Find(update.Message.Text)

			if curr.Title != "" {
				// Print curr task
				msg.Text = tm.Show(curr)
				Send(msg)

				msg.Text = Actions
				msg.ReplyMarkup = changeKeyboard
			} else {
				msg.Text = UnknownCommand
				msg.ReplyMarkup = mainKeyboard
			}
		}

		sl.Save(TaskFile(update.Message.Chat.ID))
		Send(msg)
	}
}

func main() {
	Auth()
	go NotificationOn(true)
	UpdateHandler()
}
