package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
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
	db       *botDB
	building int
}

func (cmd *setLeftStateCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	if upd.Message.IsCommand() {
		buildings := cmd.getBuildings(upd)

		if len(buildings) == 0 {
			text := fmt.Sprintf("Сначала добавьте хотя бы одну квартиру.")
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
			return msg, false
		}

		if len(buildings) > 1 {
			return cmd.askBuilding(upd, buildings)
		}

		cmd.building = buildings[0]
		return cmd.askState(upd)
	}

	if cmd.building < 0 {
		n, err := strconv.ParseInt(upd.Message.Text, 10, 32)
		if err != nil || n < 0 || n > buildingCount {
			text := fmt.Sprintf("Неверный подъезд.")
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
			return msg, false
		}

		cmd.building = int(n)
		return cmd.askState(upd)
	}

	n, err := strconv.ParseInt(upd.Message.Text, 10, 32)
	if err != nil || n < 0 || n > liftCount {
		text := fmt.Sprintf("Неверное количество.")
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		return msg, false
	}

	cmd.db.setLiftState(upd.Message.Chat.ID, cmd.building, int(n))

	text := fmt.Sprintf("Сохранено.")
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	return msg, false
}

func (cmd *setLeftStateCommand) getBuildings(upd *tgbotapi.Update) []int {
	flats := cmd.db.getUserFlats(upd.Message.Chat.ID)

	var buildings []int
	for _, flat := range flats {
		building := getBuilding(flat)

		contains := false
		for _, b := range buildings {
			if b == building {
				contains = true
				break
			}
		}
		if !contains {
			buildings = append(buildings, building)
		}
	}

	sort.Ints(buildings[:])

	return buildings
}

func (cmd *setLeftStateCommand) askBuilding(upd *tgbotapi.Update, buildings []int) (tgbotapi.MessageConfig, bool) {
	markup := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow())
	markup.OneTimeKeyboard = true
	for _, b := range buildings {
		btn := tgbotapi.NewKeyboardButton(strconv.Itoa(b))
		markup.Keyboard[0] = append(markup.Keyboard[0], btn)
	}

	text := fmt.Sprintf("Выберите подъезд.")
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ReplyMarkup = markup
	return msg, true
}

func (cmd *setLeftStateCommand) askState(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	text := fmt.Sprintf("Сколько лифтов работает в подъезде %d?", cmd.building)
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ReplyMarkup = liftStateMarkup
	return msg, true
}

type setLeftStateCommandCreator struct {
	db *botDB
}

func (cc *setLeftStateCommandCreator) BotCommand() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "set_lift_state",
		Description: "установить состояние лифтов",
	}
}

func (cc *setLeftStateCommandCreator) Create() Command {
	return &setLeftStateCommand{db: cc.db, building: -1}
}

func newSetLiftStateCommandCreator(db *botDB) CommandCreator {
	return &setLeftStateCommandCreator{db: db}
}
