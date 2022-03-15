package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type rmFlatCommand struct {
}

func (cmd *rmFlatCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	if upd.Message.IsCommand() {
		flats := []int{5, 683, 792}
		text := fmt.Sprintf("Выберите квартиру для удаления.")

		markup := tgbotapi.NewReplyKeyboard()
		markup.OneTimeKeyboard = true
		for _, flat := range flats {
			row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(strconv.Itoa(flat)))
			markup.Keyboard = append(markup.Keyboard, row)
		}

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		msg.ReplyMarkup = markup
		return msg, true
	}

	n, err := strconv.ParseInt(upd.Message.Text, 10, 32)
	flat := int(n)
	if err != nil || flat < minFlat || flat > maxFlat {
		text := fmt.Sprintf("Неверный номер.")
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		return msg, false
	}

	text := fmt.Sprintf("Квартира %d удалена.", flat)
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

type rmFlatCommandCreator struct {
}

func (cc *rmFlatCommandCreator) Text() string {
	return "rm_flat"
}

func (cc *rmFlatCommandCreator) Create() Command {
	return &rmFlatCommand{}
}

func newRmFlatCommandCreator() CommandCreator {
	return &rmFlatCommandCreator{}
}
