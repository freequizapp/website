package main

import( 
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
	"github.com/rdhmdhl/quizai/routes"
)

func main(){
	fmt.Println("hello world")

	// load env file
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Error("no .env file found -- continuing anyways...")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowedOrigins, ","),
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-auth-token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	routes.RegisterAllRoutes(r)
	r.Run("0.0.0.0:8080")
}
