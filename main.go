package main

import (
	"net/http"
	"os"
	"qrgen/service/database"
	"qrgen/service/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// init function invokes before the main function
func init() {
	// Load function loads the env file in the directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load the env files")
		return
	}
}

func main() {
	// connect to database
	dbConn, err := database.ConnectToMongoDB()
	if err != nil {
		log.Panic().Err(err).Msg("could not connect to the database")
		return
	}

	// creates an Engine instance
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// registers the routers in the routes package
	routes.RegisterAPI(router, dbConn)

	// listen and serve on 0.0.0.0:8080 ("localhost:8080")~
	router.Run(os.Getenv("SEVERICE_PORT"))
}
