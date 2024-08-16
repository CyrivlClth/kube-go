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
	err = d.load("examples/values.yaml")
	assert.NoError(t, err)
	conf := model.AppConfig{}
	assert.NoError(t, db.Where("file_name=?", "values.yaml").First(&conf).Error)
	assert.EqualValues(t, "values.yaml", conf.FileName)
}
