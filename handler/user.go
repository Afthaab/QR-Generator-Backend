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
	_, _, err = service.QrCodeGen(studentData.Id, studentData.Name)
	if err != nil {
		log.Error().Err(err).Msg("could not generate QR code")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate QR code",
		})
		return
	}

	// Step 7: Send an email with the generated QR code
	err = service.SendEmail(studentData.Email, studentData.Name+".png")
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

	// Return the student as JSON
	c.JSON(http.StatusOK, gin.H{
		"student detail": studentDetail,
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
