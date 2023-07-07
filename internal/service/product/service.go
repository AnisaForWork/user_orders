package product

import (
	"context"
	"io"
	"time"

	"github.com/AnisaForWork/user_orders/internal/config"
	"github.com/signintech/gopdf"
)

// Repository used to call db level logic
type Repository interface {
}

// AService struct implements auth service functionality
type PService struct {
	Repo           Repository
	pathToCheckDir string
}

// Order used to parse into JSON response
type Product struct {
	Barcode  string
	Name     string
	Descr    string
	Cost     int
	Created  *time.Time
	FileName string
}

func NewService(repo Repository, cfg config.Product) *PService {

	s := &PService{
		Repo:           repo,
		pathToCheckDir: cfg.PathToCheckDir,
	}
	return s
}

// Create new product for user with given login
// returns error if somethoing went wrong
func (s *PService) Create(ctx context.Context, pr Product, login string) error {
	return nil
}

// UserProducts  returns user products
// prodsPerPage - number of products on page
// page - next page with products
func (s *PService) UserProducts(ctx context.Context, page, prodsPerPage int, login string) ([]Product, error) {
	return nil, nil
}

// UserProduct returns user product
func (s *PService) UserProduct(ctx context.Context, barcode string, login string) (*Product, error) {
	return nil, nil
}

// Delete delets user product
func (s *PService) Delete(ctx context.Context, barcode string, login string) error {
	return nil
}

// GenCheck creates check for product? saves in configured directory and send PDF to user
func (s *PService) GenCheck(ctx context.Context, barcode string, login string) (gopdf.GoPdf, error) {
	return gopdf.GoPdf{}, nil
}

func (s *PService) UserProductCheck(ctx context.Context, filename string, login string) (io.Reader, error) {
	return nil, nil
}
