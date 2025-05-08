package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"Name_IQ_Finder/internal/entity"
)

func RegisterRoutes(router *gin.Engine, useCase entity.PersonUseCase) {
	handler := NewPersonHandler(useCase)

	v1 := router.Group("/api/v1")
	{
		persons := v1.Group("/persons")
		{
			persons.POST("", handler.Create)
			persons.GET("", handler.GetAll)
			persons.GET("/:id", handler.GetByID)
			persons.PUT("/:id", handler.Update)
			persons.DELETE("/:id", handler.Delete)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
