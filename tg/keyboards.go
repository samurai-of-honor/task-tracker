package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "task-manager/tg/localization"
)

var mainKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(ShowTasks),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(AddTask),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(ChangeTasks),
	),
)

var showKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(UncompletedTasks),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(AllTasks),
		tg.NewKeyboardButton(OverdueTasks),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(Return),
	),
)

var changeKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(MarkTask),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(ChangeTask),
		tg.NewKeyboardButton(DeleteTask),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(Return),
	),
)

var optionKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(ChangeTitle),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(ChangeDesc),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton(ChangeDLIne),
	),
)
