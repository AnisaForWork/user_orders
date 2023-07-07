package product

import (
	"context"
	"io"
	"regexp"

	"github.com/AnisaForWork/user_orders/internal/handler/error/mapper"
	"github.com/AnisaForWork/user_orders/internal/service/product"
	"github.com/signintech/gopdf"

	"github.com/gin-gonic/gin"
)

// Service used to call auth service level auth
type Service interface {
	Create(ctx context.Context, pr product.Product, login string) error
	UserProducts(ctx context.Context, page, prodsPerPage int, login string) ([]product.Product, error)
	Delete(ctx context.Context, login string, barcode string) error
	UserProduct(ctx context.Context, barcode string, login string) (*product.Product, error)
	GenCheck(ctx context.Context, barcode string, login string) (gopdf.GoPdf, error)
	UserProductCheck(ctx context.Context, filename string, login string) (io.Reader, error)
}

type Router struct {
	service      Service
	errMapper    mapper.ErrorMapper
	barcodeRegex *regexp.Regexp
}

func NewRouter(service Service) *Router {
	mapping := mapper.NewErrorMapper(
		mapper.ErrorMap{},
	)

	barcodeRegex := regexp.MustCompile(`^[0-9]{10}$`)

	router := &Router{
		service:      service,
		errMapper:    mapping,
		barcodeRegex: barcodeRegex,
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
