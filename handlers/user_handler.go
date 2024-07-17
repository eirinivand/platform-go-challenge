package handlers

import (
	"encoding/json"
	"errors"
	"favourites/database"
	"favourites/models"
	"favourites/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"time"
)

type UserHandler struct {
	service database.UserService
}

func NewUserHandler(service database.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(ctx *gin.Context) {
	users, err := h.service.GetAll(nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByUsername(ctx *gin.Context) {
	username, _ := ctx.Params.Get("username")
	user, err := h.service.GetByUsername(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) Add(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		fmt.Println(err)
	}
	var result *models.User
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}

	err = h.service.Create(ctx, result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *UserHandler) AddAll(ctx *gin.Context) {
	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var result []*models.User
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}

	err = h.service.CreateAll(ctx, result)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	ctx.Status(http.StatusCreated)

}

var jwtKey = []byte(utils.JwtSecret)

func (h *UserHandler) Login(ctx *gin.Context) {

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := ctx.Params.Get("username")
	existingUser, err := h.service.GetByUsername(ctx, username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Username,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	ctx.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"success": "user logged in"})
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := ctx.Params.Get("username")
	user, err := h.service.GetByUsername(ctx, username)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate password hash"})
		return
	}

	err = h.service.Create(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate password hash"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "user created"})
}
func (h *UserHandler) LogOut(ctx *gin.Context) {

	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"success": "user logged out"})
}

func (h *UserHandler) Home(ctx *gin.Context) {

	cookie, err := ctx.Cookie("token")

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, err := utils.ParseToken(cookie)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "favourites page", "role": claims.Role})
}
