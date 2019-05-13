package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pqnguyen/simple-chatapp/types"
	"log"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("sqlite3", "chatapp.db")
	if err != nil {
		log.Fatalf("error while connect to database: %s", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&Message{})
}

type Message struct {
	gorm.Model
	Receiver int
	Sender   int
	Message  string
	Status   types.MessageStatus
}
