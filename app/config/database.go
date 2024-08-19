package config

import (
	"github.com/CyrivlClth/kube-go/app/query"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("kubego.db"))
	if err != nil {
		return nil, err
	}
	query.SetDefault(db)
	return db, err
}
