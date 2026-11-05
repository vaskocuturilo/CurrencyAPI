package controller

import (
	"CurrencyAPI/models"
	"CurrencyAPI/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CurrencyController struct {
	service service.ICurrencyService
}

func NewUserController(service service.ICurrencyService) *CurrencyController {
	return &CurrencyController{service: service}
}

func (c *CurrencyController) Get(ctx *gin.Context) {
	currency := ctx.Param("currency")

	result, err := c.service.Get(ctx.Request.Context(), currency)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			errorResponse(ctx, http.StatusNotFound, "The currency is not found")
		default:
			errorResponse(ctx, http.StatusInternalServerError, models.InternalError)
			return
		}
	}

	ctx.JSON(http.StatusOK, models.CurrencyResponse{
		Date:  result.Date,
		Base:  result.Base,
		Quote: result.Quote,
		Rate:  result.Rate,
	})
}
func (c *CurrencyController) GetAll(ctx *gin.Context) {
	result, err := c.service.GetAll(ctx.Request.Context())

	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, models.InternalError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
func (c *CurrencyController) SaveAll(ctx *gin.Context) {
	var currencies []models.CurrencyResponse

	err := c.service.SaveAll(ctx.Request.Context(), currencies)

	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, models.InternalError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Currencies list saved successfully"})
}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
