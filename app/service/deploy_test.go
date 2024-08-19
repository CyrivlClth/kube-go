package service

import (
	"testing"

	"github.com/CyrivlClth/kube-go/app/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeploy_Load_UpsertOK(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	d := NewDeploy(db.Debug())
	err = d.Load("examples/values.yaml")
	assert.NoError(t, err)
	conf := model.EnvConfig{}
	assert.NoError(t, db.Where("file_name=?", "values.yaml").First(&conf).Error)
	assert.EqualValues(t, "values.yaml", conf.FileName)
	err = d.Load("examples/values.yaml")
	assert.NoError(t, err)
}

func TestDeploy_DeployApp_SameNameUpsert(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	d := NewDeploy(db.Debug())
	assert.NoError(t, d.Load("examples/values.yaml"))
	var c int64
	envName := "values.yaml"
	assert.NoError(t, db.Model(&model.EnvConfig{}).Where("file_name=?", envName).Count(&c).Error)
	assert.EqualValues(t, 1, c)
	app := model.AppConfig{
		Name: "gateway-service",
		AppBaseConfig: model.AppBaseConfig{
			MaxCPUCount: 2,
			MaxMemoryGB: 2,
			Description: "网关服务",
			PreCmd:      []string{"tini", "java"},
			Args:        []string{},
			PostCmd:     []string{"-jar", "./app.jar"},
			NodeSelector: map[string]any{
				"resources.type/base": "true",
			},
			Replicas: 3,
		},
	}
	assert.NoError(t, d.AddApp(&app))
	assert.NoError(t, db.Model(&model.AppConfig{}).Where("name=?", app.Name).Count(&c).Error)
	assert.EqualValues(t, 1, c)
	assert.Error(t, d.AddApp(&app))
	assert.NoError(t, db.Model(&model.AppConfig{}).Where("name=?", app.Name).Count(&c).Error)
	assert.EqualValues(t, 1, c)
	dp := model.AppDeploy{
		AppName: app.Name,
		EnvName: envName,
		Image:   "test",
		Tag:     "v1",
	}
	assert.NoError(t, d.DeployApp(&dp))
	actual := model.AppDeploy{}
	assert.NoError(t, d.db.Where("app_name=?", dp.AppName).Where("env_name=?", dp.EnvName).First(&actual).Error)
	assert.EqualValues(t, dp, actual)
	assert.NoError(t, d.db.Model(&model.AppDeploy{}).Where("app_name=?", dp.AppName).Where("env_name=?", dp.EnvName).Count(&c).Error)
	assert.EqualValues(t, 1, c)

	dp = model.AppDeploy{
		AppName: app.Name,
		EnvName: envName,
		Image:   "test",
		Tag:     "v2",
	}
	assert.NoError(t, d.DeployApp(&dp))
	actual = model.AppDeploy{}
	assert.NoError(t, d.db.Where("app_name=?", dp.AppName).Where("env_name=?", dp.EnvName).First(&actual).Error)
	assert.EqualValues(t, dp, actual)
	assert.NoError(t, d.db.Model(&model.AppDeploy{}).Where("app_name=?", dp.AppName).Where("env_name=?", dp.EnvName).Count(&c).Error)
	assert.EqualValues(t, 1, c)

	dp = model.AppDeploy{
		AppName: app.Name,
		EnvName: envName,
		Image:   "dev",
		Tag:     "v1",
	}
	assert.NoError(t, d.DeployApp(&dp))
	actual = model.AppDeploy{}
	assert.NoError(t, d.db.Where("app_name=?", dp.AppName).Where("env_name=?", dp.EnvName).First(&actual).Error)
	assert.EqualValues(t, dp, actual)
	assert.NoError(t, d.db.Model(&model.AppDeploy{}).Where("app_name=?", dp.AppName).Where("env_name=?", dp.EnvName).Count(&c).Error)
	assert.EqualValues(t, 1, c)
}

func TestDeploy_ExportEnv(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	d := NewDeploy(db.Debug())
	err = d.Load("examples/values.yaml")
	assert.NoError(t, err)
	conf := model.EnvConfig{}
	assert.NoError(t, db.Where("file_name=?", "values.yaml").First(&conf).Error)
	assert.EqualValues(t, "values.yaml", conf.FileName)
	app := model.AppConfig{
		Name: "gateway-service",
		AppBaseConfig: model.AppBaseConfig{
			MaxCPUCount: 2,
			MaxMemoryGB: 2,
			Description: "网关服务",
			PreCmd:      []string{"tini", "java"},
			Args:        []string{},
			PostCmd:     []string{"-jar", "./app.jar"},
			NodeSelector: map[string]any{
				"resources.type/base": "true",
			},
			Replicas: 3,
		},
	}
	assert.NoError(t, d.AddApp(&app))
	dp := model.AppDeploy{
		AppName: app.Name,
		EnvName: conf.FileName,
		Image:   "nginx-alpine",
		Tag:     "v2",
	}
	assert.NoError(t, d.DeployApp(&dp))
	assert.NoError(t, d.ExportEnv(conf.FileName))
}
