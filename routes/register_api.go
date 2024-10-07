package routes

import (
	"net/http"
	"qrgen/service/handler"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAPI(router *gin.Engine, dbConn *mongo.Database) {
	handlerFunc := handler.NewHandler(dbConn)

	// health check router to check if the server is working fine
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// user routes
	// router.POST("/user/create", handlerFunc.CreateUser)

	// admin routes
	// router.POST("/admin/register", handlerFunc.RegisterAdmin)
	router.POST("/teacher/signin", handlerFunc.TeacherSignIn)
}
