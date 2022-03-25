package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)


func (U *UserHandler) CreateUser(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Please send proper request body"})
	}

	ctx := context.Background()

	var userExist User
	searchErr := U.UserCollection.FindOne(ctx,bson.M{"name": user.Name}).Decode(&userExist)
	if searchErr != nil && searchErr != mongo.ErrNoDocuments{ 
		fmt.Println(searchErr)
		return err  
	}else if searchErr == nil && userExist.Name == user.Name{
		return c.JSON(http.StatusAlreadyReported, "User exists, log in to continue using")
	}

	user.Password, _ = HashPassword(user.Password)
	_, err = U.UserCollection.InsertOne(ctx,user)
	if err != nil{
		fmt.Println("CreateUser Insertion Failed %+v\n",err)
		return c.JSON(http.StatusInternalServerError, "DB insertion failed")
	}

	token, err := CreateToken(user)
	if err != nil {
		fmt.Println("Error creating token: %+v\n", err)
		return err
	}

	c.Response().Header().Set("x-auth-token", token)	
	return c.JSON(http.StatusCreated, user)
}

func (U *UserHandler) SigninUser(c echo.Context) error {
	ctx := context.Background()

	var userReq, user User 
	err := c.Bind(&userReq)
	if err != nil {
		fmt.Println("Error binding request body %+v\n", err)
		return err 
	}
	
	err = U.UserCollection.FindOne(ctx, bson.M{"name": userReq.Name}).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments{
		fmt.Println("Error binding request body %+v\n", err)
		return err 
	}else if err == mongo.ErrNoDocuments{
		fmt.Printf("User not found %+v\n", err)
		return err
	}

	err = ComparePassword(user.Password, userReq.Password)
	if err != nil {
		fmt.Println("Incorrect password %+v\n", err)
		return err 
	}

	token, err := CreateToken(user)
	if err != nil {
		fmt.Println("Error creating token: %+v\n", err)
		return err
	}

	user.Password = ""

	c.Response().Header().Set("x-auth-token", token)
	return c.JSON(http.StatusAccepted, user)
}

func HashPassword(password string) (string, error) {
	pwd := []byte(password)
	res, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	return string(res), nil
}

func ComparePassword(hashedPassword string, password string) error {
	pwd := []byte(password)
	hash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(hash, pwd)
	return err
}

func CreateToken(user User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user"] = user.Name
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte("secret"))
	if err != nil {
		fmt.Printf("Error signing jwt token %+v\n", err)
		return "",err
	}

	return token, nil
}