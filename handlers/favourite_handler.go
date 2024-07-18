package handlers

import (
	"encoding/json"
	"favourites/database"
	"favourites/models"
	"favourites/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type FavouriteHandler struct {
	service database.FavouriteService
}

func NewFavouriteHandler(service database.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{service: service}
}

func (h *FavouriteHandler) GetAll(ctx *gin.Context) {
	favourites, err := h.service.GetAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if favourites == nil {
		// will return "null" if empty, with this "trick" we return "[]" json.
		favourites = make([]models.Favourite, 0)
	}

	ctx.JSON(http.StatusOK, favourites)
}

func (h *FavouriteHandler) Get(ctx *gin.Context) {
	id, _ := ctx.Params.Get("id")

	m, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == utils.ErrorNotFound {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, m)
}

func (h *FavouriteHandler) Add(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}

	var result *models.Favourite
	err = json.Unmarshal(byteValue, &result)

	result.Role = ctx.GetString("role")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.service.Create(ctx, result)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}
