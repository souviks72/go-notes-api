package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	userHandler := handlers.UserHandler{UserCollection: notesDB.Collection("users")}

	config := middleware.JWTConfig{
        SigningKey: []byte("secret"),
    }
	
	e.POST("/note", notesHandler.CreateNote, middleware.JWTWithConfig(config))
	e.GET("/notes", notesHandler.GetNote, middleware.JWTWithConfig(config))
	e.GET("/note/:id", notesHandler.GetNoteById, middleware.JWTWithConfig(config))
	e.DELETE("/note/:id", notesHandler.DeleteNote, middleware.JWTWithConfig(config))
	e.PATCH("/note/:id", notesHandler.EditNote, middleware.JWTWithConfig(config))

	e.POST("/user/signup", userHandler.CreateUser)
	e.POST("/user/signin", userHandler.SigninUser)
	e.Logger.Fatal(e.Start(":8080"))
}