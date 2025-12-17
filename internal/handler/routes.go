package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, handler *DeviceHandler) {
    // Swagger route
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    api := r.Group("/api/v1")
    {
        api.POST("/devices", handler.CreateDevice)
        api.GET("/devices/:id", handler.GetDevice)
        api.GET("/devices", handler.ListDevices)
        api.PUT("/devices/:id", handler.UpdateDevice) // Assuming PUT for full/partial if flexible
        api.PATCH("/devices/:id", handler.UpdateDevice) // Mapping PATCH too
        api.DELETE("/devices/:id", handler.DeleteDevice)
    }
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
}
