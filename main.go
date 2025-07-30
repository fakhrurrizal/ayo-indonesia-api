package main

import (
	router "ayo-indonesia-api/app/routers"
	"ayo-indonesia-api/config"
	"log"
	"net/http"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

// @title Ayo Indonesia API
// @description API documentation by Muhammad Fakhrur Rizal

// @securityDefinitions.apikey JwtToken
// @in header
// @name Authorization

func main() {
	app := router.Init()

	config.Database()

	app.Static("/assets", "./assets")

	addr := "127.0.0.1:" + config.LoadConfig().Port
	server := &http.Server{
		Addr:    addr,
		Handler: app,
	}

	log.Printf("Server: %s", config.LoadConfig().BaseUrl)
	log.Printf("Documentation: %s/docs", config.LoadConfig().BaseUrl)

	graceful.ListenAndServe(server, 5*time.Second)
}
