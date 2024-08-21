package handler

import (
	"github.com/CyrivlClth/kube-go/app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(e *gin.Engine, db *gorm.DB) {
	apiG := e.Group("/api")
	{
		svc := service.NewDeploy(db)
		h := NewDeploy(svc)
		apiG.GET("/-/reload", WrapHandle(h.Reload))
		apiG.POST("/app-config", WrapHandle(h.AddApp))
		apiG.GET("/app-config", WrapHandle(h.ListApp))
		apiG.POST("/app-deploy", WrapHandle(h.DeployApp))
	}
}
