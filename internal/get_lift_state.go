package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
	"strconv"
	"time"
)

type getLeftStateCommand struct {
	db       *botDB
	building int
}

func (cmd *getLeftStateCommand) Exec(upd *tgbotapi.Update) (tgbotapi.MessageConfig, bool) {
	if upd.Message.IsCommand() {
		buildings := cmd.getBuildings(upd)

		if len(buildings) == 0 {
			buildings = []int{1, 2, 3, 4}
		}

		if len(buildings) > 1 {
			return cmd.askBuilding(upd, buildings)
		}

		cmd.building = buildings[0]
	}

	if cmd.building < 0 {
		n, err := strconv.ParseInt(upd.Message.Text, 10, 32)
		if err != nil || n < 0 || n > buildingCount {
			text := fmt.Sprintf("Неверный подъезд.")
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
			return msg, false
		}

		cmd.building = int(n)
	}

	n, at := cmd.db.getLiftState(cmd.building)
	if n < 0 {
		text := fmt.Sprintf("Нет информации.")
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
		return msg, false
	}

	text := fmt.Sprintf("В подъезде %d ", cmd.building)
	switch n {
	case 0:
		text += "не работает <b>ни один</b> лифт."
	case 1:
		text += "работает <b>один</b> лифт."
	case 2:
		text += "работают <b>два</b> лифта."
	case 3:
		text += "работают <b>все</b> лифты."
	}

	text += fmt.Sprintf(" Обновлено %s.", formatDate(at))

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	return msg, false
}

func (cmd *getLeftStateCommand) getBuildings(upd *tgbotapi.Update) []int {
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

func (cmd *getLeftStateCommand) askBuilding(upd *tgbotapi.Update, buildings []int) (tgbotapi.MessageConfig, bool) {
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

type getLeftStateCommandCreator struct {
	db *botDB
}

func (cc *getLeftStateCommandCreator) BotCommand() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "get_lift_state",
		Description: "узнать состояние лифтов",
	}
}

func (cc *getLeftStateCommandCreator) Create() Command {
	return &getLeftStateCommand{db: cc.db, building: -1}
}

func newGetLiftStateCommandCreator(db *botDB) CommandCreator {
	return &getLeftStateCommandCreator{db: db}
}

func formatDate(date time.Time) string {
	today := time.Now()

	if today.YearDay() == date.YearDay() && today.Year() == date.Year() {
		return "сегодня в " + date.Format("15:04")
	}

	yesterday := today.Add(-24 * time.Hour)
	if yesterday.YearDay() == date.YearDay() && yesterday.Year() == date.Year() {
		return "вчера в " + date.Format("15:04")
	}

	var str = date.Format("_2")

	switch date.Month() {
	case 1:
		str += " января"
	case 2:
		str += " февраля"
	case 3:
		str += " марта"
	case 4:
		str += " апреля"
	case 5:
		str += " мая"
	case 6:
		str += " июня"
	case 7:
		str += " июля"
	case 8:
		str += " августа"
	case 9:
		str += " сентября"
	case 10:
		str += " октября"
	case 11:
		str += " ноября"
	case 12:
		str += " декабря"
	}

	if today.Year() != date.Year() {
		str += " " + date.Format("2006")
	}

	return str
}
