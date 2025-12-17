package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, handler *DeviceHandler) {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    api := r.Group("/api/v1")
    {
        api.POST("/devices", handler.CreateDevice)
        api.GET("/devices/:id", handler.GetDevice)
        api.GET("/devices", handler.ListDevices)
        api.PUT("/devices/:id", handler.UpdateDevice)
        api.PATCH("/devices/:id", handler.UpdateDevice)
        api.DELETE("/devices/:id", handler.DeleteDevice)
    }
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
}
