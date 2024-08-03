package model

import "go.mongodb.org/mongo-driver/mongo"

type MongoDB struct {
	Client *mongo.Client
	Database    *mongo.Database
	Collection *mongo.Collection
}

type Result struct {
	DuplicatedKey  []interface{}
	InsertedKey    []interface{} 
}