package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var notes = []Note{
	{1, "Galaxy Note 20 Ultra","Excellent Phone"},
	{2, "Redmi Note 10 Pro","Good Budget Phone"},
}

func createNote(c echo.Context) error {
	note := new(Note)
	err := c.Bind(&note)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Please send proper request body"})
	}

	notes = append(notes, *note)
	return c.JSON(http.StatusCreated, note)
}

func getNote(c echo.Context) error {
	return  c.JSON(http.StatusOK, notes)
}

func getNoteById(c echo.Context) error {
	id := c.Param("id")

	if id == ""{
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Missing id in request path"})
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Missing id in request path"})
	}

	for _, n := range notes {
		if n.ID == idInt {
			return c.JSON(http.StatusOK, n)
		}
	}

	return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Id not found"})
}

func deleteNote(c echo.Context) error {
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,map[string]string{"msg": "Invalid id in request path"})
	}

	for i, n := range notes {
		if n.ID == id {
			notes = append(notes[0:i],notes[i+1:]...)
			return c.JSON(http.StatusOK, notes)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"msg": "ID not found"})
}

func editNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,map[string]string{"msg": "Invalid id in request path"})
	}

	note := &Note{}
	err = c.Bind(note)
	if err != nil {
		return c.JSON(http.StatusBadRequest,map[string]string{"msg": "Invalid request body"})
	}
	
	for i, n := range notes{
		if n.ID == id {
			
			if note.Content !="" {
				n.Content = note.Content
			}

			if note.Title !="" {
				n.Title = note.Title
			}
			fmt.Println(n)
			notes[i] = n 
			return c.JSON(http.StatusOK, n)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"msg": "Id no found"})
}

func main() {
	e := echo.New()
	e.POST("/note", createNote)
	e.GET("/notes", getNote)
	e.GET("/note/:id", getNoteById)
	e.DELETE("/note/:id", deleteNote)
	e.PATCH("/note/:id", editNote)
	e.Logger.Fatal(e.Start(":8080"))
}