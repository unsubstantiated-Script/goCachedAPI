package routes

import (
	"goCachedAPI/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, productHandler *handlers.ProductHandler) {
	router.POST("/product", productHandler.CreateOrUpdate)
	router.GET("/product/:id", productHandler.GetByID)
	router.PUT("/product/:id", productHandler.Update)
	router.DELETE("/product/:id", productHandler.Delete)
	router.POST("/product/:id/invalidate", productHandler.Invalidate)
	router.GET("/recent", productHandler.GetRecentProducts)
}
