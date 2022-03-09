package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Note struct {        
	Title        string `json:"title" bson:"title"`
	Content      string `json:"content" bson:"content"`
	DateCreated  string `json:"date_created,omitempty" bson:"date_created,omitempty"`
	DateModified string `json:"date_modified,omitempty" bson:"date_modified,omitempty"`
}

type NotesHandler struct {
	NotesCollection *mongo.Collection
}