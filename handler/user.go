package handler

import (
	"context"
	"net/http"
	"qrgen/service/model"
	"qrgen/service/service"
	"qrgen/service/utilities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type hanlderLayer struct {
	dbconn *mongo.Database
}

func NewHandler(dbConn *mongo.Database) hanlderLayer {
	return hanlderLayer{
		dbconn: dbConn,
	}
}

func (h *hanlderLayer) TeacherSignIn(c *gin.Context) {
	teacherSignInData := model.TeacherSignIn{}
	err := c.BindJSON(&teacherSignInData)
	if err != nil {
		log.Error().Err(err).Msg("could not bind the request body with the struct")
		utilities.BindJsonErrorResponse(c) // returning the error
		return
	}

	collection := h.dbconn.Collection("admin")

	filter := bson.M{"email": teacherSignInData.Email}

	var teacherData model.TeacherData

	err = collection.FindOne(context.Background(), filter).Decode(&teacherData)
	if err != nil {
		log.Error().Err(err).Msg("could not find the email address in the database")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email Address not found",
		})
		return
	}

	if teacherData.Password != teacherSignInData.Password {
		log.Error().Err(err).Msg("password does not match")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong Password",
		})
		return
	}

	// err = bcrypt.CompareHashAndPassword([]byte(adminData.Password), []byte(teacherSignInData.Password))
	// if err != nil {
	// 	log.Error().Err(err).Msg("password did not match")
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "incorrect password",
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"message": "success",
		"name":    teacherData.Name,
	})

}

func (h *hanlderLayer) StudentRegister(c *gin.Context) {
	studentData := model.Student{}
	err := c.BindJSON(&studentData)
	if err != nil {
		log.Error().Err(err).Msg("could not bind the request body with the struct")
		utilities.BindJsonErrorResponse(c) // returning the error
		return
	}

	collection := h.dbconn.Collection("class10")
	result, err := collection.InsertOne(context.Background(), studentData)
	if err != nil {
		log.Error().Err(err).Msg("could not bind the request body with the struct")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not create the user",
		})
		return
	}

	// get the userid
	studentData.Id = result.InsertedID.(primitive.ObjectID).Hex()

	// generate the qrCode scanner for the user id
	_, _, err = service.QrCodeGen(studentData.Id, studentData.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	err = service.SendEmail("afthab606@gmail.com", studentData.Name+".png")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": studentData.Id,
		"name":   studentData.Name,
	})

}

// func (h *hanlderLayer) RegisterAdmin(c *gin.Context) {
// 	adminData := model.TeacherData{}

// 	err := c.BindJSON(&adminData)
// 	if err != nil {
// 		log.Error().Err(err).Msg("could not bind the request body with the struct")
// 		utilities.BindJsonErrorResponse(c) // returning the error
// 		return
// 	}

// 	collection := h.dbconn.Collection("admin")

// 	filter := bson.M{"email": adminData.Email}

// 	err = collection.FindOne(context.Background(), filter).Decode(&adminData)
// 	if err == nil {
// 		log.Error().Err(err).Msg("email already exists in the database")
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "email address already exists",
// 		})
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminData.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		log.Error().Err(err).Msg("could not hash the password")
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "invalid request",
// 		})
// 		return
// 	}

// 	adminData.Password = string(hashedPassword)

// 	_, err = collection.InsertOne(context.Background(), adminData)
// 	if err != nil {
// 		log.Error().Err(err).Msg("could not bind the request body with the struct")
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "could not create the user",
// 		})
// 		return
// 	}

// 	c.JSON(200, gin.H{
// 		"message": "successfully registered",
// 	})
// }

// func (h *hanlderLayer) CreateUser(c *gin.Context) {
// 	userData := model.User{}
// 	// marshall the request
// 	err := c.BindJSON(&userData)
// 	if err != nil {
// 		log.Error().Err(err).Msg("could not bind the request body with the struct") // log the error
// 		utilities.BindJsonErrorResponse(c)                                          // returning the error
// 		return
// 	}

// 	// create the user
// 	collection := h.dbconn.Collection("class10")
// 	result, err := collection.InsertOne(context.Background(), userData)
// 	if err != nil {
// 		log.Error().Err(err).Msg("could not bind the request body with the struct")
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "could not create the user",
// 		})
// 		return
// 	}

// 	// get the userid
// 	userData.Id = result.InsertedID.(primitive.ObjectID).Hex()

// 	// generate the qrCode scanner for the user id
// 	_, _, err = service.QrCodeGen(userData.Id, userData.Name)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 		return
// 	}

// 	err = service.SendEmail("afthab606@gmail.com", userData.Name+".png")

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"userId": userData.Id,
// 		"name":   userData.Name,
// 	})

// }
