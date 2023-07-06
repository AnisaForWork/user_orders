package auth

import (
	"context"
	"net/http"

	"github.com/AnisaForWork/user_orders/internal/handler/error/mapper"
	"github.com/AnisaForWork/user_orders/internal/service/auth"

	"github.com/gin-gonic/gin"
)

// Service used to call auth service level auth
type Service interface {
	SingUp(ctx context.Context, newUser *auth.User) (err error)
	SingIn(ctx context.Context, newUser auth.Auth) (string, error)
}

type Router struct {
	service   Service
	errMapper mapper.ErrorMapper
}

func NewRouter(service Service) *Router {
	mapping := mapper.NewErrorMapper(
		mapper.ErrorMap{
			auth.ErrUserAlredyExist: {StatusCode: http.StatusConflict, Msg: "User with such credentials already exist"},
			auth.ErrNoSuchUser:      {StatusCode: http.StatusBadRequest, Msg: "Wrong login or password"},
			auth.ErrUsrDelted:       {StatusCode: http.StatusForbidden, Msg: "Forbidden, create account and log in"},
		},
	)

	router := &Router{
		service:   service,
		errMapper: mapping,
	}
	return router
}

func (a *Router) InitRoutes() *gin.Engine {
	r := gin.New()
	r.POST("/reg", a.singUp)
	r.POST("/auth", a.singIn)
	return r
}
