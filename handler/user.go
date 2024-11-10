package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"qrgen/service/model"
	"qrgen/service/service"
	"qrgen/service/utilities"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type hanlderLayer struct {
	dbconn *mongo.Database
}

func NewHandler(dbConn *mongo.Database) hanlderLayer {
	return hanlderLayer{
		dbconn: dbConn,
	}
}

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

func (h *hanlderLayer) StudentRegister(c *gin.Context) {
	studentData := model.Student{}

	// Step 1: Bind JSON body to struct
	if err := c.BindJSON(&studentData); err != nil {
		log.Error().Err(err).Msg("could not bind the request body to the student struct")
		utilities.BindJsonErrorResponse(c) // return the error response
		return
	}

	// Step 2: Check if the email is already in use
	collection := h.dbconn.Collection("class10")
	count, err := collection.CountDocuments(context.Background(), bson.M{"email": studentData.Email})
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

	studentData.Attandance = make(map[string]string)

	// Step 4: Insert the student data into MongoDB
	result, err := collection.InsertOne(context.Background(), studentData)
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
	studentData.Id = result.InsertedID.(primitive.ObjectID).Hex()

	// Step 6: Generate a QR code for the user ID
	_, _, err = service.QrCodeGen(studentData.Id, studentData.Id)
	if err != nil {
		log.Error().Err(err).Msg("could not generate QR code")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate QR code",
		})
		return
	}

	// Step 7: Send an email with the generated QR code
	err = service.SendEmail(studentData.Email, "QR_Codes/"+studentData.Id+".png")
	if err != nil {
		log.Error().Err(err).Msg("could not send the email with QR code")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not send email",
		})
		return
	}

	// Step 8: Send success response
	c.JSON(http.StatusOK, gin.H{
		"userId": studentData.Id,
		"name":   studentData.Name,
	})
}

func (h *hanlderLayer) ViewAllStudents(c *gin.Context) {
	// Step 1: Define the slice to store student details
	var studentDetails []model.Student

	// Step 2: Initialize MongoDB connection and collection
	collection := h.dbconn.Collection("class10")

	// Step 3: Create a find options (optional, you can customize projection, sorting, etc.)
	findOptions := options.Find()

	// Step 4: Fetch the documents
	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve the data")
		c.JSON(400, gin.H{
			"error": "could not find the data",
		})
		return
	}
	defer cursor.Close(context.Background())

	// Step 5: Iterate over the cursor and decode documents into the studentDetails slice
	for cursor.Next(context.Background()) {
		var student model.Student
		if err := cursor.Decode(&student); err != nil {
			log.Error().Err(err).Msg("could not retrieve the data")
			c.JSON(400, gin.H{
				"error": "could not find the data",
			})
			return
		}
		studentDetails = append(studentDetails, student)
	}

	// Step 6: Handle any cursor errors
	if err := cursor.Err(); err != nil {
		log.Error().Err(err).Msg("could not retrieve the data")
		c.JSON(400, gin.H{
			"error": "could not find the data",
		})
		return
	}

	c.JSON(200, gin.H{
		"student details": studentDetails,
	})

}

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
	collection := h.dbconn.Collection("class10")

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
	imagePath := fmt.Sprintf("./QR_Codes/%s.png", studentID)
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

func (h *hanlderLayer) RegisterAttendance(c *gin.Context) {
	registerData := model.RegisterAttendance{}

	// Step 1: Bind JSON body to struct
	if err := c.BindJSON(&registerData); err != nil {
		log.Error().Err(err).Msg("could not bind the request body to the object struct")
		utilities.BindJsonErrorResponse(c) // return the error response
		return
	}

	// Convert the string studentID to MongoDB ObjectID (if it's ObjectID)
	objID, err := primitive.ObjectIDFromHex(registerData.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid student ID"})
		return
	}

	var studentDetail model.Student

	// Step 2: Initialize MongoDB connection and collection
	collection := h.dbconn.Collection("class10")

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

	if studentDetail.Attandance[time.Now().Format("02-01-2006")] == "Present" {
		c.JSON(400, gin.H{
			"error": "attendance already registred",
		})
		return
	}

	presentDate := time.Now().Format("02-01-2006") // Use current date as an example
	status := "Present"

	// Update the document: add new date to Attandance map
	update := bson.M{
		"$set": bson.M{
			"attandance." + presentDate: status, // Use dot notation to set the new date in the map
		},
	}

	// Perform the update
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		// Handle other potential errors
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error updating the attandance"})
		return
	}

	if result.MatchedCount < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No documents matched the filter",
		})
		return
	}

	// Step 2: Initialize MongoDB connection and collection
	collection = h.dbconn.Collection("attandance")

	// Find document with the current date
	filter1 := bson.M{"date": presentDate}
	update1 := bson.M{"$addToSet": bson.M{"Attendees": studentDetail.Name}}

	// Try to update the document
	updateResult, err := collection.UpdateOne(context.Background(), filter1, update1)
	if err != nil {
		// Handle other potential errors
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error updating the attandance"})
		return
	}

	// Check if a document was modified
	if updateResult.MatchedCount == 0 {
		// Document not found, so insert a new one
		newDoc := model.TimeSheet{
			Date: presentDate,
			Attendees: []string{
				studentDetail.Name,
			},
		}

		_, err := collection.InsertOne(context.Background(), newDoc)
		if err != nil {
			// Check if it's a MongoDB duplicate key error (E11000)
			if mongo.IsDuplicateKeyError(err) {
				log.Error().Err(err).Msg("duplicate key error: date or another field might be unique")
				c.JSON(http.StatusConflict, gin.H{
					"error": "duplicate record exists for the timesheet",
				})
				return
			}

			log.Error().Err(err).Msg("could not insert attendacne data into the collection")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not create the attendance record",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"success": "successfully registered the attendance",
	})
}

func (h *hanlderLayer) ViewAllAttendacne(c *gin.Context) {
	// Get the studentId from the URL
	date := c.Param("date")
	attendanceData := model.TimeSheet{}

	collection := h.dbconn.Collection("attandance")
	err := collection.FindOne(context.TODO(), bson.M{"date": date}).Decode(&attendanceData)
	if err == mongo.ErrNoDocuments {
		// If no document found
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found for the give date"})
		return
	} else if err != nil {
		// Handle other potential errors
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error finding the date"})
		return
	}

	c.JSON(200, gin.H{
		"attendees": attendanceData.Attendees,
	})
}
