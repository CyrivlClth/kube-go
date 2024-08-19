package handler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/CyrivlClth/kube-go/app/model"
	"github.com/CyrivlClth/kube-go/app/service"
)

func NewDeploy(svc *service.Deploy) *Deploy {
	root, err := filepath.Abs(filepath.Join("./envconfig"))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(root, os.ModePerm)
	if err != nil {
		panic(err)
	}
	deploy := &Deploy{
		root: root,
		svc:  svc,
	}
	return deploy
}

type Deploy struct {
	root string
	svc  *service.Deploy
}

func (d Deploy) Reload(c *Context) error {
	dir, err := os.ReadDir(d.root)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(dir))
	for _, f := range dir {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".yaml") {
			err := d.svc.Load(filepath.Join(d.root, f.Name()))
			if err != nil {
				return err
			}
			files = append(files, f.Name())
		}
	}
	c.JSON(files)
	return nil
}

func (d Deploy) AddApp(c *Context) error {
	var r model.AppConfig
	if err := c.ShouldBindJSON(&r); err != nil {
		return err
	}
	if err := d.svc.AddApp(&r); err != nil {
		return err
	}
	c.JSON(r)
	return nil
}

func (d Deploy) ListApp(c *Context) error {
	apps, err := d.svc.ListApp("")
	if err != nil {
		return err
	}
	c.JSON(apps)
	return nil
}

func (d Deploy) DeployApp(c *Context) error {
	var r model.AppDeploy
	if err := c.ShouldBindJSON(&r); err != nil {
		return err
	}
	if err := d.svc.DeployApp(&r); err != nil {
		return err
	}
	c.JSON(r)
	return nil
}
