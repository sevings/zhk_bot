package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type rmFlatCommand struct {
	db *botDB
}

func (cmd *rmFlatCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	if upd.Message.IsCommand() {
		flats := cmd.db.getUserFlats(upd.Message.Chat.ID)
		if len(flats) == 0 {
			text := fmt.Sprintf("Нет добавленных квартир.")
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
			return msg, false
		}

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

	cmd.db.removeUserFlat(upd.Message.Chat.ID, flat)

	text := fmt.Sprintf("Квартира %d удалена.", flat)
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

type rmFlatCommandCreator struct {
	db *botDB
}

func (cc *rmFlatCommandCreator) Text() string {
	return "rm_flat"
}

func (cc *rmFlatCommandCreator) Create() Command {
	return &rmFlatCommand{db: cc.db}
}

func newRmFlatCommandCreator(db *botDB) CommandCreator {
	return &rmFlatCommandCreator{db: db}
}
