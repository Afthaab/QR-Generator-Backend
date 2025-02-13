package routes

import (
	"net/http"
	"qrgen/service/auth"
	"qrgen/service/handler"
	"qrgen/service/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAPI(router *gin.Engine, dbConn *mongo.Database, auth auth.Auth) {
	handlerFunc := handler.NewHandler(dbConn, auth)

	// health check router to check if the server is working fine
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	m := middleware.NewMiddleware(auth)

	// user routes
	// router.POST("/user/create", handlerFunc.CreateUser)

	// admin routes
	// router.POST("/admin/register", handlerFunc.RegisterAdmin)
	router.POST("/teacher/signin", handlerFunc.TeacherSignIn)
	router.POST("/teacher/student/register", handlerFunc.StudentRegister)
	router.GET("/student/view/all", handlerFunc.ViewAllStudents)
	router.GET("/student/view", m.Authenticate(handlerFunc.ViewStudent))
	router.POST("/teacher/signup", handlerFunc.TeacherSignup)
	router.POST("/register/attendance", handlerFunc.RegisterAttendance)
	router.GET("/view/all/attendance/:date", handlerFunc.ViewAllAttendacne)
	router.POST("/student/sign/in", handlerFunc.StudentSignIn)
}
