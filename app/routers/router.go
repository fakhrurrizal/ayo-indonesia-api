package router

import (
	"ayo-indonesia-api/app/controllers"
	"ayo-indonesia-api/app/middlewares"
	"ayo-indonesia-api/config"
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

	router.GET("/", controllers.Index)

	v1 := router.Group("/v1", middlewares.StripHTMLMiddleware())
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/signin", controllers.SignIn)
			auth.GET("/user", middlewares.Auth(), controllers.GetSignInUser)
			auth.PUT("/user-profile", middlewares.Auth(), controllers.UpdateUserProfileByID)
		}

		// Trip routes
		trip := v1.Group("/trip")
		{
			trip.POST("", middlewares.Auth(), controllers.CreateTrip)
			trip.GET("/:id", controllers.GetTripByID)
			trip.GET("", controllers.GetTrips)
			trip.DELETE("/:id", middlewares.Auth(), controllers.DeleteTripByID)
			trip.PUT("/:id", middlewares.Auth(), controllers.UpdateTripByID)
		}

		// File routes
		file := v1.Group("/file", middlewares.Auth())
		{
			file.POST("", controllers.UploadFile)
			file.GET("", controllers.GetFile)
		}

		// Destination Type routes
		destinationType := v1.Group("/destination-type")
		{
			destinationType.POST("", middlewares.Auth(), controllers.CreateDestinationType)
			destinationType.GET("/:id", controllers.GetDestinationTypeByID)
			destinationType.GET("", controllers.GetDestinationTypes)
			destinationType.DELETE("/:id", middlewares.Auth(), controllers.DeleteDestinationTypeByID)
			destinationType.PUT("/:id", middlewares.Auth(), controllers.UpdateDestinationTypeByID)
		}
	}

	log.Println("Server initialized...")

	return router
}
