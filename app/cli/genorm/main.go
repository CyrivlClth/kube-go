package main

import (
	"github.com/CyrivlClth/kube-go/app/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

//go:generate go run .
func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	gormdb, _ := gorm.Open(sqlite.Open(":memory:"))
	g.UseDB(gormdb)
	g.ApplyBasic(model.AppConfig{}, model.AppDeploy{}, model.EnvConfig{})
	g.Execute()
}
