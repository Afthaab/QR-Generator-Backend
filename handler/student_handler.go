package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"qrgen/service/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *hanlderLayer) ViewStudent(c *gin.Context) {
	// Get the studentId from the URL
	studentID := c.Param("studentId")

	// Convert the string studentID to MongoDB ObjectID (if it's ObjectID)
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid student ID"})
		return
	}

	var studentDetail model.Student

	// Step 2: Initialize MongoDB connection and collection
	collection := h.dbconn.Collection("students")

	// Find the student in the MongoDB collection
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&studentDetail)
	if err == mongo.ErrNoDocuments {
		// If no document found
		c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		return
	} else if err != nil {
		// Handle other potential errors
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error finding student"})
		return
	}

	// Load the image from local folder and encode it in base64
	imagePath := fmt.Sprintf("./qr_codes/%s.png", studentID)
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading student image"})
		return
	}

	// Encode image to base64
	studentDetail.Image = base64.StdEncoding.EncodeToString(imageData)

	// Return the student as JSON
	c.JSON(http.StatusOK, gin.H{
		"student detail": studentDetail,
	})
}
