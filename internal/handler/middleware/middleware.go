package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AnisaForWork/user_orders/internal/handler/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var KeyUserID = "userID"

const (
	AuthHeaderKey   = "Authorization"
	AuthTypeBearer  = "bearer"
	authTokenFields = 2
	InvalidUser     = int64(-1)
)

// Service used to call auth service level auth
type Service interface {
	ValidateToken(header string) error
	UserLogin(tokenStr string) (string, error)
}

type Middleware struct {
	service Service
}

func NewMiddleware(srv Service) *Middleware {
	m := &Middleware{
		service: srv,
	}

	return m
}

type correlationIDType int

const (
	requestIDKey correlationIDType = iota
)

// WithRqID returns a context which knows its request ID
func WithRqID(ctx context.Context, rqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, rqID)
}

// Logger middleware adds logger with id to the request context and logs request execution time
func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		ctx := WithRqID(c.Request.Context(), uuid.NewString())
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			logrus.WithContext(ctx).
				WithFields(logrus.Fields{
					"request":       c.Request.RequestURI,
					"executionTime": time.Since(t),
					"status":        c.Writer.Status(),
					"clientip":      c.ClientIP(),
				}).Info("Request finished in time")
		}()

		c.Next()
	}
}

// Recovery middleware adds recovery strategy(logs error)
func (m *Middleware) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.String(http.StatusInternalServerError, "internal error")
		//log.WithCtx(c.Request.Context()).Error("Error during execution", zap.String("place", "request"), zap.Any("panic", recovered), zap.String("path", c.FullPath()))
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// Authentication middleware if request provided valid token extracts userID and adds it to context
func (m *Middleware) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if c.IsAborted() {
				c.Set(KeyUserID, InvalidUser)
			}
		}()

		log := logrus.WithContext(c.Request.Context())

		accessToken, err := ExtractToken(c)
		log.WithFields(logrus.Fields{
			"token": accessToken,
		}).Info("Processing access token")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.CreateJSONResult("Error", err.Error()))
			return
		}

		if err := m.service.ValidateToken(accessToken); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.CreateJSONResult("Error", err.Error()))
			return
		}

		login, err := m.service.UserLogin(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.CreateJSONResult("Error", err.Error()))
			return
		}

		c.Set(KeyUserID, login)

		c.Next()
	}
}

func ExtractToken(c *gin.Context) (string, error) {
	token := c.GetHeader(AuthHeaderKey)

	fields := strings.Fields(token)
	if len(fields) != authTokenFields {
		return "", fmt.Errorf("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != AuthTypeBearer {
		return "", fmt.Errorf("unsupported authorization type %s", authorizationType)
	}

	return fields[1], nil
}
