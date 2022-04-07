package handlers

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type Note struct {        
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"`
	DateCreated  string             `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateModified string             `json:"date_modified,omitempty" bson:"date_modified,omitempty"`
	Username     string             `json:"user_name" bson:"user_name"`
}

type NotesHandler struct {
	NotesCollection *mongo.Collection
}

type User struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password,omitempty" bson:"password"`
}

type UserHandler struct {
	UserCollection *mongo.Collection
}

type ClaimsStruct struct{
	UserName string `json:"username"`
	Authorised bool `json:"authorised"`
	jwt.StandardClaims
}