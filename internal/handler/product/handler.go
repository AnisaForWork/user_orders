package product

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/AnisaForWork/user_orders/internal/handler/error/validator"
	"github.com/AnisaForWork/user_orders/internal/handler/middleware"
	"github.com/AnisaForWork/user_orders/internal/handler/response"
	"github.com/AnisaForWork/user_orders/internal/service/product"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Created used to parse requests body with info about new Order
type Created struct {
	Barcode string `json:"barcode"   binding:"required,len=10,numeric" minimum:"10" maximum:"10" default:"1234567890"`
	Name    string `json:"name"  binding:"required,min=10,max=60,startsnotwith= ,endsnotwith= " minimum:"10" maximum:"60" default:"product name"`
	Descr   string `json:"desc" binding:"required,min=10,max=1000,startsnotwith= ,endsnotwith= " default:"Description"`
	Cost    int    `json:"cost" binding:"required,min=1" minimum:"10" default:"100"`
}

// Product model used to parse into JSON response
type Product struct {
	Barcode  string     `json:"barcode,omitempty"`
	Name     string     `json:"name,omitempty"`
	Descr    string     `json:"desc,omitempty"`
	Cost     int        `json:"cost,omitempty"`
	Created  *time.Time `json:"created,omitempty"`
	FileName string     `json:"filaname,omitempty"`
}

// @Summary      Create new product
// @Description  user provides products barcode, name, description and cost
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        product      body  product.Created true "barcode,name,desc,cost"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      409  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/ [post]
func (p *Router) create(c *gin.Context) {
	log := logrus.WithContext(c.Request.Context())

	login := c.GetString(middleware.KeyUserID)

	var prodt Created
	if err := c.ShouldBindJSON(&prodt); err != nil {
		c.JSON(http.StatusBadRequest, validator.ProcessValidatorError(err))
		return
	}

	srvProdt := product.Product{
		Barcode: prodt.Barcode,
		Name:    prodt.Name,
		Descr:   prodt.Descr,
		Cost:    prodt.Cost,
	}

	err := p.service.Create(c.Request.Context(), srvProdt, login)
	if err != nil {
		log.WithFields(logrus.Fields{
			"handler":   "product",
			"func":      "create",
			"userLogin": login,
		}).Error("Error during creating product")

		errInf := p.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	c.JSON(http.StatusCreated, response.CreateJSONResult("Product created", []string{login, prodt.Barcode}))
}

// @Summary      Returns all users products with pagnation
// @Description  configured to return set number of products(barcodes+name+cost) per request, user sends configures how many products per page be seen and offset, only owner of product can view it
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param   	 p query     int    true "Next page to retrieve" minimum(1)    maximum(50000)
// @Param   	 n query     int    true "Number of products info per page" minimum(1)    maximum(100)
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      403  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/all [get]
func (p *Router) allUserProducts(c *gin.Context) {

	log := logrus.WithContext(c.Request.Context())

	login := c.GetString(middleware.KeyUserID)

	page, err := strconv.Atoi(c.Query("p"))
	if err != nil || (page < 0 || page > 50000) {
		c.JSON(http.StatusBadRequest, validator.ErrorMsg("p", "shouldd be between 1 and 500000"))
		return
	}

	prodsPerPage, err := strconv.Atoi(c.Query("n"))
	if err != nil || (prodsPerPage < 1 || prodsPerPage > 0) {
		c.JSON(http.StatusBadRequest, validator.ErrorMsg("n", "shouldd be between 1 and 100"))
		return
	}

	prods, err := p.service.UserProducts(c.Request.Context(), page, prodsPerPage, login)
	if err != nil {
		log.WithFields(logrus.Fields{
			"handler":    "prodyct",
			"func":       "allUserProducts",
			"userLogin":  login,
			"pages":      p,
			"numPerPage": prodsPerPage,
		}).Error("Error retrieving user products")

		errInf := p.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	c.JSON(http.StatusOK, response.CreateJSONResult("Products", prods))
}

// @Summary      Returns user product full info
// @Description  returns full info about product user chose to view, only owner of product can view it
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param 		 barcode   path      string true  "Product barcode"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      404  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/{id} [get]
func (p *Router) userProduct(c *gin.Context) {
	log := logrus.WithContext(c.Request.Context())

	login := c.GetString(middleware.KeyUserID)

	barcode := c.Query("barcode")
	if !p.barcodeRegex.MatchString(barcode) {
		c.JSON(http.StatusBadRequest, validator.ErrorMsg("barcode", "should consist of ten numeric numbers"))
		return
	}

	prod, err := p.service.UserProduct(c.Request.Context(), barcode, login)
	if err != nil {
		log.WithFields(logrus.Fields{
			"handler":   "product",
			"func":      "create",
			"userLogin": login,
			"barcode":   barcode,
		}).Error("Error retrieving product")

		errInf := p.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	c.JSON(http.StatusCreated, response.CreateJSONResult("Product", prod))
}

// @Summary      Delete user product
// @Description  archives product but don't delete it from storage, only owner of product can do it
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param 		 barcode   path      string true  "Product barcode"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      404  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/{id} [get]
func (p *Router) delete(c *gin.Context) {
	log := logrus.WithContext(c.Request.Context())

	login := c.GetString(middleware.KeyUserID)

	barcode := c.Query("barcode")
	if !p.barcodeRegex.MatchString(barcode) {
		c.JSON(http.StatusBadRequest, validator.ErrorMsg("barcode", "should consist of ten numeric numbers"))
		return
	}

	err := p.service.Delete(c.Request.Context(), barcode, login)
	if err != nil {
		log.WithFields(logrus.Fields{
			"handler":   "product",
			"func":      "create",
			"userLogin": login,
			"barcode":   barcode,
		}).Error("Error deleting product")

		errInf := p.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	c.JSON(http.StatusCreated, response.CreateJSONResult("Succesfull", "Product deleted"))
}

// @Summary      Generates check
// @Description  using info about given product generates PDF check using special PDF template
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param 		 barcode   path      string true  "Product barcode"
// @Produce  	 application/pdf
// @Success 	 200 {file} PdfFile
// @Failure      400  {object}  response.JSONResult
// @Failure      403  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/{id}/check [get]
func (p *Router) genCheck(c *gin.Context) {
	log := logrus.WithContext(c.Request.Context())

	login := c.GetString(middleware.KeyUserID)

	barcode := c.Query("barcode")
	if !p.barcodeRegex.MatchString(barcode) {
		c.JSON(http.StatusBadRequest, validator.ErrorMsg("barcode", "should consist of ten numeric numbers"))
		return
	}

	pdrWriter, err := p.service.GenCheck(c.Request.Context(), barcode, login)
	if err != nil {
		log.WithFields(logrus.Fields{
			"handler":   "product",
			"func":      "create",
			"userLogin": login,
			"barcode":   barcode,
		}).WithError(err).Error("Error creating check for product")

		errInf := p.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	c.Header("Content-type", "application/pdf")
	c.Status(http.StatusOK)
	err = pdrWriter.Write(c.Writer)
	pdrWriter.Close()
	if err != nil {
		log.WithError(err).Warn("Could not send check")
	}
}

// @Summary      Returns check of user product
// @Description  sends user PDF check generated previously using product info, only owner of product can do it
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param 		 checkName   path   string true  "Check file name"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      403  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/check/{checkName} [get]
func (p *Router) productCheck(c *gin.Context) {
	log := logrus.WithContext(c.Request.Context())

	login := c.GetString(middleware.KeyUserID)

	barcode := c.Query("barcode")
	if !p.barcodeRegex.MatchString(barcode) {
		c.JSON(http.StatusBadRequest, validator.ErrorMsg("barcode", "should consist of ten numeric numbers"))
		return
	}

	f, err := p.service.UserProductCheck(c.Request.Context(), barcode, login)
	if err != nil {
		log.WithFields(logrus.Fields{
			"handler":   "product",
			"func":      "create",
			"userLogin": login,
			"barcode":   barcode,
		}).Error("Error sending check file")

		errInf := p.errMapper.MapError(err)
		c.JSON(errInf.StatusCode,
			response.CreateJSONResult("Error", errInf.Msg))

		return
	}

	c.Header("Content-type", "application/pdf")
	c.Status(http.StatusOK)
	//Stream to response
	if _, err := io.Copy(c.Writer, f); err != nil {
		c.JSON(http.StatusBadRequest, validator.ErrorMsg("barcode", "should consist of ten numeric numbers"))
	}
}
