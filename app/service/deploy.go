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
	err := db.AutoMigrate(&model.AppConfig{}, &model.EnvConfig{})
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
		EnvConfig model.EnvConfig `yaml:"envConfig"`
	}{}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return err
	}
	_, conf.EnvConfig.FileName = filepath.Split(root)
	err = d.db.Create(&(conf.EnvConfig)).Error

	return err
}
