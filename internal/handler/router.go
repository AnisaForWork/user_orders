package handler

import (
	docs "github.com/AnisaForWork/user_orders/api/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Service holds all needed services
type Service interface {
}

// @title           User orders service API
// @version         1.0
// @description     User orders service API using swagger 2.0.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// NewRouter returns instane with paths of requests mapped on their handlers
func NewRouter(service Service, isRelease bool) *gin.Engine {
	if isRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	docs.SwaggerInfo.BasePath = "/"
	return router
}
