package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/AnisaForWork/user_orders/internal/config"
	"github.com/AnisaForWork/user_orders/internal/repository/mysql"

	"github.com/signintech/gopdf"
)

// Repository used to call db level logic
type Repository interface {
	Create(ctx context.Context, pr mysql.Product, login string) error
	UserProducts(ctx context.Context, amount int, offset int, usrID string) ([]mysql.Product, error)
	UserProduct(ctx context.Context, barcode string, login string) (mysql.Product, error)
	Delete(ctx context.Context, barcode string, login string) error
	ProfuctInfoForCheck(ctx context.Context, barcode string, login string) (mysql.Product, error)
	UpdateCheckInfo(ctx context.Context, filename string, barcode string, login string) error
	CheckOwnership(ctx context.Context, filename string, login string) error
}

var (
	ErrNotOwner  = errors.New("user is not the owner of product")
	ErrNotExists = errors.New("product dosn't exist")
)

// AService struct implements auth service functionality
type PService struct {
	Repo           Repository
	PathToCheckDir string
	PathToTemplate string
	TemplateName   string
	FontName       string
	FontFileName   string
	TimeFormat     string
	TemplateW      float64
	TemplateH      float64
}

// Order used to parse into JSON response
type Product struct {
	Barcode  string
	Name     string
	Descr    string
	Cost     int
	Created  time.Time
	FileName string
}

func NewService(repo Repository, cfg config.Product) *PService {

	s := &PService{
		Repo:           repo,
		PathToCheckDir: cfg.PathToCheckDir,
		PathToTemplate: cfg.PathToTemplate,
		TemplateName:   cfg.TemplateName,
		FontName:       cfg.FontName,
		FontFileName:   cfg.FontFileName,
		TimeFormat:     cfg.TimeFormat,
		TemplateH:      cfg.TemplateH,
		TemplateW:      cfg.TemplateW,
	}
	return s
}

// Create new product for user with given login
// returns error if somethoing went wrong
func (s *PService) Create(ctx context.Context, pr Product, login string) error {
	dbModel := mysql.Product{
		Barcode: pr.Barcode,
		Name:    pr.Name,
		Desc:    pr.Descr,
		Cost:    pr.Cost,
	}
	return s.Repo.Create(ctx, dbModel, login)
}

// UserProducts  returns user products
// prodsPerPage - number of products on page
// page - next page with products
func (s *PService) UserProducts(ctx context.Context, page, prodsPerPage int, login string) ([]Product, error) {
	products, err := s.Repo.UserProducts(ctx, prodsPerPage, (page-1)*prodsPerPage, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	if len(products) < 1 {
		return nil, nil
	}

	resProds := make([]Product, len(products))
	for i := 0; i < len(products); i++ {
		resProds[i] = Product{
			Barcode: products[i].Barcode,
			Name:    products[i].Name,
			Cost:    products[i].Cost,
		}
	}

	return resProds, nil

}

// UserProduct returns user product
func (s *PService) UserProduct(ctx context.Context, barcode string, login string) (*Product, error) {
	dbModel, err := s.Repo.UserProduct(ctx, barcode, login)
	if err != nil {
		return nil, err
	}

	res := &Product{
		Barcode:  dbModel.Barcode,
		Name:     dbModel.Name,
		Descr:    dbModel.Desc,
		Cost:     dbModel.Cost,
		Created:  dbModel.Created,
		FileName: dbModel.FileName,
	}
	return res, nil
}

// Delete delets user product
func (s *PService) Delete(ctx context.Context, barcode string, login string) error {
	return s.Repo.Delete(ctx, barcode, login)
}

// GenCheck creates check for product? saves in configured directory and send PDF to user
func (s *PService) GenCheck(ctx context.Context, barcode string, login string) (gopdf.GoPdf, error) {
	prod, err := s.Repo.ProfuctInfoForCheck(ctx, barcode, login)

	if err != nil {
		return gopdf.GoPdf{}, err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: s.TemplateW, H: s.TemplateW}})
	pdf.AddPage()

	path := filepath.Join(s.PathToTemplate, s.TemplateName)
	// import template file
	tpl := pdf.ImportPage(path, 1, "/MediaBox") // 1 is the page
	// Draw pdf onto page
	pdf.UseImportedTemplate(tpl, 0, 0, s.TemplateW, s.TemplateW) // Template structure, x coordinate, y coordinate, width, height

	path = filepath.Join(s.PathToTemplate, s.FontFileName)
	err = pdf.AddTTFFont(s.FontName, path)
	if err != nil {
		return gopdf.GoPdf{}, err
	}

	err = pdf.SetFont(s.FontName, "", 10)
	if err != nil {
		return gopdf.GoPdf{}, err
	}
	pdf.SetXY(21, 36)
	pdf.Text(prod.Barcode) // y coordinate specification

	pdf.SetFontSize(8)
	pdf.SetXY(21, 75)
	pdf.Text(prod.Name) // y coordinate specification

	pdf.SetFontSize(10)
	pdf.SetXY(161, 116)
	pdf.Text(strconv.Itoa(prod.Cost)) // y coordinate specification

	fileName := fmt.Sprintf("doc_%s_%s.pdf", prod.Barcode, time.Now().Format(s.TimeFormat))

	path = filepath.Join(s.PathToCheckDir, fileName)
	err = s.Repo.UpdateCheckInfo(ctx, fileName, barcode, login)
	if err != nil {
		return gopdf.GoPdf{}, err
	}
	pdf.WritePdf(path)
	return pdf, nil
}

func (s *PService) UserProductCheck(ctx context.Context, fileName string, login string) (io.Reader, error) {
	err := s.Repo.CheckOwnership(ctx, fileName, login)
	if err != nil {
		return nil, ErrNotOwner
	}

	path := filepath.Join(s.PathToCheckDir, fileName)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
