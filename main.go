package main

import (
	_ "time"
	"reflect"
	"./log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type TestTable struct {
	ID			uint `gorm:"primary_key"`
}

func (TestTable) TableName() string {
	return "t_test"
}

func checkTableExist(db *gorm.DB) {
	tablePtr := &TestTable{}
	log.Debug("Type: %s", reflect.TypeOf(*tablePtr))

	if db.HasTable(tablePtr) {
		log.Info("Has table")
		return
	} else {
		log.Info("Does not have table")
	}
	// does not have table
	// db.Set("ENGINE=InnoDB", "DEFAULT CHARSET=utf8").CreateTable(tablePtr)
	// log.Info("Create table")
}

func testReflect(v interface{}) {
	log.Debug("Type of v: %s", reflect.TypeOf(v));
}

func main() {
	log.Debug("Hello, gorm!")

	db, err := gorm.Open("mysql", "tars:tars2015@tcp(10.0.4.11:3306)/db_gorm?charset-utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		log.Error("Error opening DB: %s", err.Error())
	} else {
		log.Info("Open DB OK!")
	}

	checkTableExist(db)
	testReflect(db)
	return
}
