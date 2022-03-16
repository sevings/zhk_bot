package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

var liftStateMarkup tgbotapi.ReplyKeyboardMarkup

func init() {
	liftStateMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("0"),
			tgbotapi.NewKeyboardButton("1"),
			tgbotapi.NewKeyboardButton("2"),
			tgbotapi.NewKeyboardButton("3"),
		))
	liftStateMarkup.OneTimeKeyboard = true
}

type setLeftStateCommand struct {
}

func (cmd *setLeftStateCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	building := getBuilding(683)

	if upd.Message.IsCommand() {
		text := fmt.Sprintf("Сколько лифтов работает в %d корпусе на данный момент?", building)
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		msg.ReplyMarkup = liftStateMarkup
		return msg, true
	}

	n, err := strconv.ParseInt(upd.Message.Text, 10, 32)
	if err != nil || n < 0 || n > liftCount {
		text := fmt.Sprintf("Неверное количество.")
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		return msg, false
	}

	text := fmt.Sprintf("Сохранено.")
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

type setLeftStateCommandCreator struct {
}

func (cc *setLeftStateCommandCreator) Text() string {
	return "set_lift_state"
}

func (cc *setLeftStateCommandCreator) Create() Command {
	return &setLeftStateCommand{}
}

func newSetLiftStateCommandCreator() CommandCreator {
	return &setLeftStateCommandCreator{}
}
