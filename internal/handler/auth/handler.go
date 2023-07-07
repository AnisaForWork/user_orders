package auth

import (
	"net/http"

	"github.com/AnisaForWork/user_orders/internal/handler/error/validator"
	"github.com/AnisaForWork/user_orders/internal/handler/response"
	"github.com/AnisaForWork/user_orders/internal/service/auth"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Information that needed to register new user, stored in request body
type Registration struct {
	Login    string `json:"login" binding:"required,min=3,max=40,startsnotwith= ,endsnotwith= " minimum:"3" maximum:"40" example:"Login123"`
	FullName string `json:"fullName" binding:"required,min=3,max=75,startsnotwith= ,endsnotwith= " minimum:"3" maximum:"75" example:"Ivanov Ivan Ivanovich"`
	Email    string `json:"email" binding:"required,email" maximum:"255"  example:"test@test.com" `
	Password string `json:"password" binding:"required,min=6,max=40"  minimum:"6" maximum:"40" example:"password"`
}

// @Summary      Register user
// @Description  singUp with credentials user : login,full name,email,password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   reg      body     auth.Registration true "login,full name,email,password"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      409  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /auth/reg [post]
func (a *Router) singUp(c *gin.Context) {
	log := logrus.WithContext(c.Request.Context())

	var reg Registration
	if err := c.ShouldBindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, validator.ProcessValidatorError(err))
		return
	}

	user := auth.User{
		Login:    reg.Login,
		FullName: reg.FullName,
		Password: reg.Password,
		Email:    reg.Email,
	}

	err := a.service.SingUp(c.Request.Context(), &user)
	if err != nil {
		log.WithFields(logrus.Fields{
			"handlers": "auth",
			"func":     "singUp",
		}).WithError(err).Error("Error during creating user")

		errInf := a.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	log.WithFields(logrus.Fields{
		"Login": reg.Login,
	}).Info("User registered")
	c.JSON(http.StatusCreated, response.CreateJSONResult("Successfully registered", nil))
}

// Auth used to parse requests body from singing in request
type Auth struct {
	Login    string `json:"login" binding:"required,min=3,max=40,startsnotwith= ,endsnotwith= " minimum:"3" maximum:"40" example:"Login123"`
	Password string `json:"password" binding:"required,min=6,max=40"  minimum:"6" maximum:"40" default:"password"`
}

// @Summary      Sing in user
// @Description  sing in user if they have given valid credentials, returns access token(JWT)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth  body     auth.Auth true "login,password"
// @Success      200  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      404  {object}  response.JSONResult
// @Failure      403  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /auth/auth [post]
func (a *Router) singIn(c *gin.Context) {
	log := logrus.WithContext(c.Request.Context())

	var au Auth
	if err := c.ShouldBindJSON(&au); err != nil {
		c.JSON(http.StatusBadRequest, validator.ProcessValidatorError(err))
		return
	}

	srvAuth := auth.Auth{
		Login:    au.Login,
		Password: au.Password,
	}

	srvTokens, err := a.service.SingIn(c.Request.Context(), srvAuth)

	if err != nil {
		log.WithFields(logrus.Fields{
			"handlers": "auth",
			"func":     "singIn",
		}).WithError(err).Error("Error during singin in user")

		errInf := a.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	log.WithFields(logrus.Fields{
		"UserToken": srvTokens,
	}).Info("User received tokens")

	c.JSON(http.StatusOK, response.CreateJSONResult("Access token", srvTokens))
}
