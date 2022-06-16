package main

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var mainKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("🖼 Show tasks"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("📌 Add task"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("🛠 Change tasks"),
	),
)

var showKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("🗃 All"),
		tg.NewKeyboardButton("🛑 Uncompleted"),
		tg.NewKeyboardButton("‼️Overdue"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("🔙 Return"),
	),
)

var changeKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("✅ Mark task"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("🔧 Change task"),
		tg.NewKeyboardButton("🗑 Delete task"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("🔙 Return"),
	),
)

var optionKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("✏️ Title"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("📝 Description"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("⏰ Deadline"),
	),
)
