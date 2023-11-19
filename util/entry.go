package util

import (
	"fmt"
	"uno/pkg/database"
	"uno/service/subscriber"
)

// Setup initializes the util
func Setup() {
	if !database.GetWriteDB().Migrator().HasTable(&subscriber.Entry{}) {
		fmt.Println("create table subscriber")
		database.GetWriteDB().Migrator().CreateTable(&subscriber.Entry{})
	}
}
