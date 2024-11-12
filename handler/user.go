package handler

import (
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
