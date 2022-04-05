package internal

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type botDB struct {
	db *gorm.DB
}

type user struct {
	ID     int64  `gorm:"not null;primaryKey;unique;index"`
	Name   string `gorm:"not null;unique"`
	flats  []userFlat
	states []liftState
}

type userFlat struct {
	UserID int64 `gorm:"not null;index"`
	Flat   int   `gorm:"not null;index"`
}

type liftState struct {
	ID        uint      `gorm:"not null;primaryKey;unique;autoIncrement"`
	UpdatedAt time.Time `gorm:"not null"`
	Building  int       `gorm:"not null"`
	Working   int       `gorm:"not null"`
	UserID    int64     `gorm:"not null"`
}

func openBotDB(source string) (*botDB, error) {
	db, err := gorm.Open(sqlite.Open(source), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&user{}, &userFlat{}, &liftState{})
	if err != nil {
		return nil, err
	}

	return &botDB{db: db}, nil
}

func checkGormError(tx *gorm.DB) {
	if tx.Error != nil {
		log.Println(tx.Error)
	}
}

func (db *botDB) addUserFlat(userID int64, name string, flat int) {
	res := db.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user{ID: userID, Name: name})
	checkGormError(res)

	res = db.db.Create(&userFlat{
		UserID: userID,
		Flat:   flat,
	})

	checkGormError(res)
}

func (db *botDB) getUserFlats(userID int64) []int {
	var flats []int
	res := db.db.Model(&userFlat{}).Where("user_id = ?", userID).Pluck("flat", &flats)
	checkGormError(res)

	return flats
}

func (db *botDB) removeUserFlat(userID int64, flat int) {
	res := db.db.Where(&userFlat{UserID: userID, Flat: flat}, "user_id", "flat").Delete(&userFlat{})

	checkGormError(res)
}

func (db *botDB) setLiftState(userID int64, building, working int) {
	res := db.db.Create(&liftState{
		UpdatedAt: time.Now(),
		Building:  building,
		Working:   working,
		UserID:    userID,
	})

	checkGormError(res)
}

func (db *botDB) getLiftState(building int) (int, time.Time) {
	state := &liftState{Working: -1}
	res := db.db.Where(&liftState{Building: building}, "building").Order("updated_at desc").First(state)
	checkGormError(res)

	return state.Working, state.UpdatedAt
}
