package main

import (
	"log"
	"mysql/config"
	"mysql/model"
	"mysql/routes"
	"mysql/utils"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			result := config.DB.Where("expires_at < ? ", time.Now()).
				Delete(&model.Session{})
			log.Printf("Session cleanup: removed %d expired/revoked sessions", result.RowsAffected)
		}
	}()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-api-key", "X-Admin-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(utils.SecurityHeaders())
	routes.SetupRoutes(r)
	port := os.Getenv("PORT")
	lc := os.Getenv("LC")
	if err := r.Run(lc + port); err != nil {
		panic(err)
	}
}
