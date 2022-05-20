package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	// uso el mail como usuario
	Id       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
	Orders   *[]Order           `bson:"orders,omitempty"`
	Mails    string             `bson:"mails,omitempty"`
}
