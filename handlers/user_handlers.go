package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)


func (U *UserHandler) CreateUser(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Please send proper request body"})
	}

	res, err :=U.UserCollection.InsertOne(context.Background(),user)
	if err != nil{
		fmt.Println("CreateUser Insertion Failed %+v\n",err)
		return echo.NewHTTPError(http.StatusInternalServerError, "DB insertion failed")
	}
	fmt.Println(res)
	return c.JSON(http.StatusCreated, user)
}