package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startCommand struct {
}

func (cmd *startCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	text := `Это бот для ЖК Фонвизинский в Москве. Справка: /help.`
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

type startCommandCreator struct {
}

func (cc *startCommandCreator) BotCommand() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "start",
		Description: "",
	}
}

func (cc *startCommandCreator) Create() Command {
	return &startCommand{}
}

func newStartCommandCreator() CommandCreator {
	return &startCommandCreator{}
}
