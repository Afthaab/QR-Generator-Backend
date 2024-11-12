package handler

import (
	"context"
	"net/http"
	"qrgen/service/model"
	"qrgen/service/utilities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *hanlderLayer) TeacherSignup(c *gin.Context) {
	teacherData := model.TeacherData{}
	// Step 1: Bind JSON body to struct
	if err := c.BindJSON(&teacherData); err != nil {
		log.Error().Err(err).Msg("could not bind the request body to the teacher struct")
		utilities.BindJsonErrorResponse(c) // return the error response
		return
	}

	// check if the passwords match
	if teacherData.Password != teacherData.Repeatpassword {
		c.JSON(400, gin.H{
			"error": "passwords do not match",
		})
		return
	}

	// Step 2: Check if the email is already in use
	collection := h.dbconn.Collection("admin")
	count, err := collection.CountDocuments(context.Background(), bson.M{"email": teacherData.Email})
	if err != nil {
		log.Error().Err(err).Msg("could not check if email exists")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error   "})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email already exists",
		})
		return
	}

	// Step 4: Insert the student data into MongoDB
	result, err := collection.InsertOne(context.Background(), teacherData)
	if err != nil {
		// Check if it's a MongoDB duplicate key error (E11000)
		if mongo.IsDuplicateKeyError(err) {
			log.Error().Err(err).Msg("duplicate key error: Email or another field might be unique")
			c.JSON(http.StatusConflict, gin.H{
				"error": "duplicate record exists for the student",
			})
			return
		}

		log.Error().Err(err).Msg("could not insert student data into the collection")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create the student record",
		})
		return
	}
	// Step 5: Get the inserted user ID
	id := result.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(200, gin.H{
		"message": "successfully registered",
		"id":      id,
	})
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

	c.JSON(200, gin.H{
		"message": "success",
		"name":    teacherData.Name,
	})

}

func (h *hanlderLayer) StudentSignIn(c *gin.Context) {
	studentSignInData := model.StudentSignIn{}
	// Step 1: Bind JSON body to struct
	if err := c.BindJSON(&studentSignInData); err != nil {
		log.Error().Err(err).Msg("could not bind the request body to the StudentSignIn struct")
		utilities.BindJsonErrorResponse(c) // return the error response
		return
	}

	collection := h.dbconn.Collection("students")

	filter := bson.M{"email": studentSignInData.Email}

	var studentData model.Student

	err := collection.FindOne(context.Background(), filter).Decode(&studentData)
	if err != nil {
		log.Error().Err(err).Msg("could not find the email address in the database")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email Address not found",
		})
		return
	}

	if studentData.Password != studentSignInData.Password {
		log.Error().Err(err).Msg("password does not match")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong Password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"student": studentData,
	})

}
