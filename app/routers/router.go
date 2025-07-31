package router

import (
	"ayo-indonesia-api/app/controllers"
	"ayo-indonesia-api/app/middlewares"
	"ayo-indonesia-api/config"
	_ "ayo-indonesia-api/docs"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c *gin.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Init() *gin.Engine {
	router := gin.New()

	router.Use(middlewares.Cors())
	router.Use(middlewares.Secure())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(static.Serve("/assets", static.LocalFile("assets", true)))
	router.GET("/", controllers.Index)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/docs", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("docs.html"))
		err := tmpl.ExecuteTemplate(c.Writer, "docs.html", gin.H{
			"BaseUrl": config.LoadConfig().BaseUrl,
			"Title":   "Api Documentation of " + config.LoadConfig().AppName,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Render error: %v", err))
			return
		}
	})

	v1 := router.Group("/v1", middlewares.StripHTMLMiddleware())
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signin", controllers.SignIn)
			auth.POST("/signup", controllers.SignUp)
			auth.GET("/user", middlewares.Auth(), controllers.GetSignInUser)
		}

		file := v1.Group("/file", middlewares.Auth())
		{
			file.POST("", middlewares.Auth(), controllers.UploadFile)
			file.GET("", controllers.GetFile)
		}
		teams := v1.Group("/teams")
		{
			teams.GET("", controllers.GetTeams)
			teams.GET("/:id", controllers.GetTeamByID)
			teams.POST("", middlewares.Auth(), controllers.CreateTeam)
			teams.PUT("/:id", middlewares.Auth(), controllers.UpdateTeam)
			teams.DELETE("/:id", middlewares.Auth(), controllers.DeleteTeam)
		}

		players := v1.Group("/players")
		{
			players.GET("", controllers.GetPlayers)
			players.GET("/:id", controllers.GetPlayerByID)
			players.POST("", middlewares.Auth(), controllers.CreatePlayer)
			players.PUT("/:id", middlewares.Auth(), controllers.UpdatePlayer)
			players.DELETE("/:id", middlewares.Auth(), controllers.DeletePlayer)
		}

		matches := v1.Group("/matches")
		{
			matches.GET("", controllers.GetMatches)
			matches.GET("/:id", controllers.GetMatchByID)
			matches.POST("", middlewares.Auth(), controllers.CreateMatch)
			matches.PUT("/:id/result", middlewares.Auth(), controllers.UpdateMatchResult)
			matches.DELETE("/:id", middlewares.Auth(), controllers.DeleteMatch)
			matches.GET("/:id/report", controllers.GetMatchReport)
		}

	}

	log.Println("Server initialized...")

	return router
}
