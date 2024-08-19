package service

import (
	"log"
	"os"
	"path/filepath"

	"github.com/CyrivlClth/kube-go/app/model"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Deploy struct {
	db *gorm.DB
}

func NewDeploy(db *gorm.DB) *Deploy {
	err := db.AutoMigrate(&model.AppConfig{}, &model.EnvConfig{}, &model.AppDeploy{})
	if err != nil {
		panic(err)
	}
	return &Deploy{
		db: db,
	}
}

func (d Deploy) Load(path string) error {
	root, _ := filepath.Abs(path)
	log.Println(root)
	b, err := os.ReadFile(root)
	if err != nil {
		return err
	}
	conf := struct {
		EnvConfig model.EnvConfig `yaml:"envConfig"`
		UserGuide map[string]any  `yaml:"userGuide"`
	}{}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return err
	}
	_, conf.EnvConfig.FileName = filepath.Split(root)
	conf.EnvConfig.UserGuide = conf.UserGuide
	err = d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&(conf.EnvConfig)).Error

	return err
}

type App struct {
	MaxCPUCount  int            `json:"maxCPUCount" yaml:"maxCPUCount" gorm:"not null"`
	MaxMemoryGB  int            `json:"maxMemoryGB" yaml:"maxMemoryGB" gorm:"not null"`
	Description  string         `json:"description" yaml:"description" gorm:"not null"`
	PreCmd       []string       `json:"preCmd" yaml:"preCmd" gorm:"not null"`
	Args         []string       `json:"args" yaml:"args" gorm:"not null"`
	PostCmd      []string       `json:"postCmd" yaml:"postCmd" gorm:"not null"`
	NodeSelector map[string]any `json:"nodeSelector" yaml:"nodeSelector" gorm:"not null"`
	Replicas     int            `json:"replicas" yaml:"replicas" gorm:"not null"`
	Name         string         `json:"name" yaml:"name"`
	Image        string         `yaml:"image"`
	Tag          string         `yaml:"tag"`
}

func (d Deploy) AddApp(app *model.AppConfig) error {
	return d.db.Create(app).Error
}

func (d Deploy) DeployApp(dp *model.AppDeploy) error {
	app := model.AppConfig{}
	err := d.db.Where("name=?", dp.AppName).First(&app).Error
	if err != nil {
		return err
	}
	env := model.EnvConfig{}
	err = d.db.Where("file_name=?", dp.EnvName).First(&env).Error
	if err != nil {
		return err
	}
	err = d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(dp).Error
	if err != nil {
		return err
	}
	return nil
}

func (d Deploy) ExportEnv(envName string) error {
	env := model.EnvConfig{}
	err := d.db.Where("file_name=?", envName).First(&env).Error
	if err != nil {
		return err
	}
	var ds []model.AppDeploy
	err = d.db.Where("env_name=?", envName).Find(&ds).Error
	if err != nil {
		return err
	}
	names := make([]string, 0, len(ds))
	for _, dp := range ds {
		names = append(names, dp.AppName)
	}
	var apps []model.AppConfig
	err = d.db.Where("name in (?)", names).Find(&apps).Error
	if err != nil {
		return err
	}
	appMap := make(map[string]model.AppConfig)
	for _, app := range apps {
		appMap[app.Name] = app
	}
	var data struct {
		EnvConfig model.EnvConfig `yaml:"envConfig"`
		UserGuide map[string]any  `yaml:"userGuide"`
		Apps      []App           `yaml:"apps"`
	}
	data.EnvConfig = env
	data.UserGuide = env.UserGuide
	for _, dp := range ds {
		c := App{}
		err = copier.Copy(&c, appMap[dp.AppName])
		if err != nil {
			return err
		}
		c.Image = dp.Image
		c.Tag = dp.Tag
		data.Apps = append(data.Apps, c)
	}
	log.Printf("%+v", data)
	f, err := os.Create(data.EnvConfig.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	e := yaml.NewEncoder(f)
	e.SetIndent(2)
	err = e.Encode(data)
	return err
}
