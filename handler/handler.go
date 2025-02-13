package handler

import (
	"qrgen/service/auth"

	"go.mongodb.org/mongo-driver/mongo"
)

type hanlderLayer struct {
	dbconn *mongo.Database
	auth   auth.Auth
}

func NewHandler(dbConn *mongo.Database, auth auth.Auth) hanlderLayer {
	return hanlderLayer{
		dbconn: dbConn,
		auth:   auth,
	}
}
