package service

import (
	"testing"

	"github.com/CyrivlClth/kube-go/app/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeploy_load(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	d := NewDeploy(db.Debug())
	err = d.Load("examples/values.yaml")
	assert.NoError(t, err)
	conf := model.EnvConfig{}
	assert.NoError(t, db.Where("file_name=?", "values.yaml").First(&conf).Error)
	assert.EqualValues(t, "values.yaml", conf.FileName)
}

func TestDeploy_exportEnv(t *testing.T) {
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
