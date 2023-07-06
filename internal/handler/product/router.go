package product

import (
	"github.com/AnisaForWork/user_orders/internal/handler/error/mapper"

	"github.com/gin-gonic/gin"
)

// Service used to call auth service level auth
type Service interface {
}

type Router struct {
	service   Service
	errMapper mapper.ErrorMapper
}

func NewRouter(service Service) *Router {
	mapping := mapper.NewErrorMapper(
		mapper.ErrorMap{},
	)

	router := &Router{
		service:   service,
		errMapper: mapping,
	}

	return router
}

func (p *Router) InitRoutes() *gin.Engine {
	r := gin.New()
	r.POST("/", p.create)
	r.GET("/all", p.allUserProducts)
	r.GET("/:id", p.userProduct)
	r.DELETE("/:id", p.delete)
	r.GET("/:id/check", p.genCheck)
	r.GET("/check/:checkName", p.productCheck)
	return r
}
