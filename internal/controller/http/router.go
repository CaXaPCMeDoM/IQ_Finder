package http

import (
	"github.com/gin-gonic/gin"

	v1 "Name_IQ_Finder/internal/controller/http/v1"
	"Name_IQ_Finder/internal/entity"
)

func NewRouter(useCase entity.PersonUseCase) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	v1.RegisterRoutes(router, useCase)

	return router
}
