package service

import (
	"log"
	"os"
	"path/filepath"

	"github.com/CyrivlClth/kube-go/app/model"
	"github.com/CyrivlClth/kube-go/app/query"
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
	return query.EnvConfig.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&(conf.EnvConfig))
}

func (d Deploy) AddApp(app *model.AppConfig) error {
	return query.AppConfig.Create(app)
}

func (d Deploy) DeployApp(dp *model.AppDeploy) error {
	_, err := query.AppConfig.
		Where(query.AppConfig.Name.Eq(dp.AppName)).
		First()
	if err != nil {
		return err
	}
	_, err = query.EnvConfig.
		Where(query.EnvConfig.FileName.Eq(dp.EnvName)).
		First()
	if err != nil {
		return err
	}
	err = query.AppDeploy.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(dp)
	if err != nil {
		return err
	}
	return nil
}

func (d Deploy) ExportEnv(envName string) error {
	env, err := query.EnvConfig.
		Where(query.EnvConfig.FileName.Eq(envName)).
		First()
	if err != nil {
		return err
	}
	type App struct {
		model.AppDeploy
		model.AppConfig
	}
	var apps []App
	err = query.AppDeploy.
		Select(query.AppDeploy.ALL, query.AppConfig.ALL).
		Where(query.AppDeploy.EnvName.Eq(envName)).
		LeftJoin(query.AppConfig, query.AppConfig.Name.EqCol(query.AppDeploy.AppName)).
		Scan(&apps)
	if err != nil {
		return err
	}
	var data struct {
		EnvConfig *model.EnvConfig `yaml:"envConfig"`
		UserGuide map[string]any   `yaml:"userGuide"`
		Apps      []App            `yaml:"apps"`
	}
	data.EnvConfig = env
	data.UserGuide = env.UserGuide
	data.Apps = apps
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

func (d Deploy) ListApp(envName string) ([]*model.AppConfig, error) {
	q := query.AppConfig
	if envName == "" {
		return q.Preload(query.AppConfig.Deploy).Find()
	}
	return q.Preload(query.AppConfig.Deploy.On(query.AppDeploy.EnvName.Eq(envName))).Find()
}
