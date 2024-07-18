// Product Api:
//
//	version: 0.1
//	title: Go Challenge
//
// Schemes: http, https
// Host:
// BasePath: /api/v1
//
//	Consumes:
//	- application/json
//
// Produces:
// - application/json
// SecurityDefinitions:
//
//	Bearer:
//	 type: apiKey
//	 name: Authorization
//	 in: header
//
// swagger:meta package main
package main

import (
	"favourites/database"
	_ "favourites/docs"
	"favourites/handlers"
	"favourites/middleware"
	"favourites/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	utils.Load()
}

func setupRouter() *gin.Engine {

	db := utils.GetDB()
	//utils.CloseClientDB()
	r := gin.Default()
	var (
		// Collections.
		chartsCollection     = db.Collection("charts")
		insightsCollection   = db.Collection("insights")
		audiencesCollection  = db.Collection("audiences")
		favouritesCollection = db.Collection("favourites")
		usersCollection      = db.Collection("users")

		// Services.
		chartService     = database.NewChartService(chartsCollection)
		insightService   = database.NewInsightService(insightsCollection)
		audienceService  = database.NewAudienceService(audiencesCollection)
		favouriteService = database.NewFavouriteService(favouritesCollection)
		assetService     = database.NewAssetService(chartService, insightService, audienceService)
		userService      = database.NewUserService(usersCollection)
	)
	appGroup := r.Group("/api/v1/")
	{

		usersGroup := appGroup.Group("/users")
		{
			userHandler := handlers.NewUserHandler(userService)
			usersGroup.POST("/login", userHandler.Login)
			usersGroup.POST("/logout", userHandler.LogOut)
			usersGroup.POST("/signup", userHandler.SignUp)
			usersGroup.GET("/profile/:username", userHandler.GetByUsername)

			favouriteGroup := usersGroup.Group("/favourites")
			{
				favouriteGroup.Use(middleware.IsAuthorized())
				favouriteHandler := handlers.NewFavouriteHandler(favouriteService)
				favouriteGroup.GET("/", favouriteHandler.GetAll)
				favouriteGroup.GET("/:id", favouriteHandler.Get)
				favouriteGroup.POST("/add", favouriteHandler.Add)

			}

		}

		adminGroup := appGroup.Group("/admin")
		{
			adminGroup.Use(middleware.IsAuthorized())
			userHandler := handlers.NewUserHandler(userService)
			adminGroup.GET("/users", userHandler.GetAll)
			adminGroup.POST("/add-users-bulk", userHandler.AddAll)

		}

		assetGroup := appGroup.Group("/assets")
		{
			assetHandler := handlers.NewAssetHandler(assetService)
			assetGroup.GET("/", assetHandler.GetAll)
		}

		insightGroup := appGroup.Group("/insights")
		{
			insightHandler := handlers.NewInsightHandler(insightService)
			insightGroup.GET("/", insightHandler.GetAll)
			insightGroup.POST("/add-bulk", insightHandler.AddAll)
		}

	}

	// Ping Health Check
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "pong"})
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in localhost:8080
	r.Run(":8080")
}
