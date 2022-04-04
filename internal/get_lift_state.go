package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
			text := fmt.Sprintf("Сначала добавьте хотя бы одну квартиру.")
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
			return msg, false
		}

		if len(buildings) > 1 {
			return cmd.askBuilding(upd, buildings)
		}

		cmd.building = buildings[0]
	}

	if cmd.building < 0 {
		n, err := strconv.ParseInt(upd.Message.Text, 10, 32)
		if err != nil || n < 0 || n > buildingCount {
			text := fmt.Sprintf("Неверный корпус.")
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

	text := fmt.Sprintf("В %d корпусе ", cmd.building)
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

	text := fmt.Sprintf("Выберите корпус.")
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	msg.ReplyMarkup = markup
	return msg, true
}

type getLeftStateCommandCreator struct {
	db *botDB
}

func (cc *getLeftStateCommandCreator) Text() string {
	return "get_lift_state"
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
	case 0:
		str += " января"
	case 1:
		str += " февраля"
	case 2:
		str += " марта"
	case 3:
		str += " апреля"
	case 4:
		str += " мая"
	case 5:
		str += " июня"
	case 6:
		str += " июля"
	case 7:
		str += " августа"
	case 8:
		str += " сентября"
	case 9:
		str += " октября"
	case 10:
		str += " ноября"
	case 11:
		str += " декабря"
	}

	if today.Year() != date.Year() {
		str += " " + date.Format("2006")
	}

	return str
}
