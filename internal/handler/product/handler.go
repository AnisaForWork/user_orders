package product

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Created used to parse requests body with info about new Order
type Created struct {
	Barcode string `json:"barcode"   binding:"required,len=10,numeric" minimum:"10" maximum:"10" default:"1234567890"`
	Name    string `json:"name"  binding:"required,min=10,max=60,startsnotwith= ,endsnotwith= " minimum:"10" maximum:"60" default:"product name"`
	Descr   string `json:"desc" binding:"required,min=10,max=1000,startsnotwith= ,endsnotwith= " default:"Description"`
	Cost    int    `json:"cost" binding:"required,min=1" minimum:"10" default:"100"`
}

// Order used to parse into JSON response
type Order struct {
	Barcode string     `json:"order,omitempty"`
	Name    string     `json:"from,omitempty"`
	Descr   string     `json:"to,omitempty"`
	Cost    int        `json:"taxi,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	File    string     `json:"rating,omitempty"`
}

// @Summary      Create new product
// @Description  user provides products barcode, name, description and cost
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        reg      body  product.Created "barcode,name,desc,cost"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      409  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/ [post]
func (p *Router) create(c *gin.Context) {

}

// @Summary      Returns all users products with pagnation
// @Description  configured to return set number of products(barcodes+name+cost) per request, user sends configures how many products per page be seen and offset, only owner of product can view it
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param   	 p query     int    true "Next page to retrieve" minimum(1)    maximum(10000)
// @Param   	 n query     int    true "Number of order info per page" minimum(1)    maximum(100)
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      403  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/all [get]
func (p *Router) allUserProducts(c *gin.Context) {

}

// @Summary      Returns user product full info
// @Description  returns full info about product user chose to view, only owner of product can view it
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param 		 id   path      string true  "Product barcode"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      404  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/{id} [get]
func (p *Router) userProduct(c *gin.Context) {

}

// @Summary      Delete user product
// @Description  archives product but don't delete it from storage, only owner of product can do it
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param 		 id   path      string true  "Product barcode"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      404  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/{id} [get]
func (p *Router) delete(c *gin.Context) {

}

// @Summary      Generates check
// @Description  using info about given product generates PDF check using special PDF template
// @Tags         product
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Security     ApiKeyAuth
// @Param 		 id   path      string true  "Product barcode"
// @Success      201  {object}  response.JSONResult
// @Failure      400  {object}  response.JSONResult
// @Failure      403  {object}  response.JSONResult
// @Failure      500  {object}  response.JSONResult
// @Router       /product/{id}/check [get]
func (p *Router) genCheck(c *gin.Context) {

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

}
