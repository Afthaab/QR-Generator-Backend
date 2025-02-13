package handler

import (
	"qrgen/service/authentication"

	"go.mongodb.org/mongo-driver/mongo"
)

type hanlderLayer struct {
	dbconn *mongo.Database
	auth   authentication.Auth
}

func NewHandler(dbConn *mongo.Database, auth authentication.Auth) hanlderLayer {
	return hanlderLayer{
		dbconn: dbConn,
		auth:   auth,
	}
}
