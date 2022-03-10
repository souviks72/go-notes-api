package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/souviks72/go-notes-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongodbUri = "mongodb://127.0.0.1:27017"

func main() {
	e := echo.New()

	mongoclient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbUri))
	if err != nil {
		e.Logger.Fatalf("Unable to connect to mongodb: %+v", err)
	}

	defer func() {
		if err = mongoclient.Disconnect(context.TODO()); err != nil{
			e.Logger.Fatalf("Unable to connect to mongodb: %+v", err)
		}
	}()

	notesDB := mongoclient.Database("go-notes")
	notesHandler := handlers.NotesHandler{NotesCollection: notesDB.Collection("notes")}
	
	e.POST("/note", notesHandler.CreateNote)
	e.GET("/notes", notesHandler.GetNote)
	e.GET("/note/:id", notesHandler.GetNoteById)
	e.DELETE("/note/:id", notesHandler.DeleteNote)
	e.PATCH("/note/:id", notesHandler.EditNote)
	e.Logger.Fatal(e.Start(":8080"))
}