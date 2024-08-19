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
		hand := NewDeploy(svc)
		apiG.GET("/-/reload", WrapHandle(hand.Reload))
		apiG.POST("/app-config", WrapHandle(hand.AddApp))
		apiG.GET("/app-config", WrapHandle(hand.ListApp))
		apiG.POST("/app-deploy", WrapHandle(hand.DeployApp))
	}
}
