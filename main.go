package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dwnwp/api-email/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	wd, _ := os.Getwd()
	if err := godotenv.Load(wd + "/.env"); err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "http://localhost:8080"},
		AllowMethods: []string{http.MethodPost, http.MethodGet},
	}))

	api := e.Group("/api")
	api.GET("/health", handlers.Healthcheck)
	api.POST("/email/create", handlers.ProducerEmail())

	fmt.Println("All registered routes:")
	data := e.Routes()
	for i := 0; i < len(data); i++ {
		fmt.Printf("Method: %s Path: %s\n", data[i].Method, data[i].Path)
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}
