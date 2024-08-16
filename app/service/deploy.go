package service

import (
	"log"
	"os"
	"path/filepath"

	"github.com/CyrivlClth/kube-go/app/model"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Deploy struct {
	db *gorm.DB
}

func NewDeploy(db *gorm.DB) *Deploy {
	err := db.AutoMigrate(&model.AppConfig{})
	if err != nil {
		panic(err)
	}
	return &Deploy{
		db: db,
	}
}

func (d Deploy) load(path string) error {
	root, _ := filepath.Abs(path)
	log.Println(root)
	b, err := os.ReadFile(root)
	if err != nil {
		return err
	}
	conf := struct {
		AppConfig model.AppConfig `yaml:"appConfig"`
	}{}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return err
	}
	_, conf.AppConfig.FileName = filepath.Split(root)
	err = d.db.Create(&(conf.AppConfig)).Error

	return err
}
