package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/44t4nk1/jwt-go/api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	router       = gin.Default()
	mySigningKey = []byte(os.Getenv("ACCESS_SECRET"))
)

//A sample use
var user = User{
	ID:       1,
	Username: "44t4nk1",
	Password: "PASSWORD",
}

//Response ...
type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

//CreateToken ...
func CreateToken(userid uint64) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

//HomePage ...
func HomePage(c *gin.Context) {
	var Response = Response{
		Error:   false,
		Message: "Home Page",
	}
	c.JSON(http.StatusOK, Response)
}

//Login ...
func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid JSON Provided")
	}
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

func isAuthorised(endpoint func(c *gin.Context)) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.GetHeader("Token") != "" {
			token, err := jwt.Parse(c.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an Error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				var Response = Response{
					Error:   true,
					Message: "Invalid Signature",
				}
				c.JSON(http.StatusUnauthorized, Response)
			}
			if token.Valid {
				endpoint(c)
			}
		} else {
			var Response = Response{
				Error:   true,
				Message: "No token provided",
			}
			c.JSON(http.StatusUnauthorized, Response)
		}
	})
}

func main() {
	router.GET("/home", isAuthorised(HomePage))
	router.POST("/login", Login)
	log.Fatal(router.Run(":8080"))
}