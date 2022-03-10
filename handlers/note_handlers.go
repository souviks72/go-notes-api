package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func (N *NotesHandler) CreateNote(c echo.Context) error {
	note := new(Note)
	err := c.Bind(&note)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Please send proper request body"})
	}

	res, err := N.NotesCollection.InsertOne(context.Background(),note)
	if err != nil{
		fmt.Println("CreateNote Insertion Failed %+v\n",err)
		return echo.NewHTTPError(http.StatusInternalServerError, "DB insertion failed")
	}
	fmt.Println(res)
	return c.JSON(http.StatusCreated, note)
}

func (N *NotesHandler) GetNote(c echo.Context) error {
	ctx := context.Background()
	mongoCursor, err := N.NotesCollection.Find(ctx,bson.M{})
	if err != nil {
		fmt.Println("CreateNote Fetch Failed %+v\n",err)
		return echo.NewHTTPError(http.StatusInternalServerError, "DB read failed")
	}

	notes := []Note{}
	defer mongoCursor.Close(ctx)
	for mongoCursor.Next(ctx){
		var n Note 
		if err = mongoCursor.Decode(&n); err != nil{
			return echo.NewHTTPError(http.StatusInternalServerError, "Decode failed")
		}
		notes = append(notes, n)
	}
	return c.JSON(http.StatusOK, notes)
}

func (N *NotesHandler) GetNoteById(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Missing id in request path"})
	}

	objId , err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get object id from hex")
	}
	res := N.NotesCollection.FindOne(context.Background(), bson.M{"_id": objId})

	var note Note
	_ = res.Decode(&note)

	return c.JSON(http.StatusBadRequest, note)
}

func (N *NotesHandler) DeleteNote(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex((c.Param("id")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Invalid id in request path"})
	}

	res := N.NotesCollection.FindOneAndDelete(context.Background(), bson.M{"_id": id})
	fmt.Println(res)

	return c.JSON(http.StatusOK, map[string]string{"msg": "Deleted"})
}

func (N *NotesHandler) EditNote(c echo.Context) error {
	var note Note
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Invalid id in request path"})
	}

	editFieldMap := map[string]interface{}{}
	err = c.Bind(&editFieldMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Invalid request body"})
	}

	editBsonMap := bson.M{}
	for k,v := range editFieldMap{
		editBsonMap[k] = v
	}
	update := bson.M{
		"$set": editBsonMap,
	}

	after := options.After
	upsert := true
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert: &upsert,
	}
	
	res := N.NotesCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": id},update, &opt)
	res.Decode(&note)

	return c.JSON(http.StatusNotFound, note)
}