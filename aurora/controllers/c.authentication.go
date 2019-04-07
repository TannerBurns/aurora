package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"

	"../models"
)

/*
AuthenticationMiddleware - middleware to authenticate using JWT,
	returns 400 for errors
*/
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1],
					func(token *jwt.Token) (interface{}, error) {
						_, ok := token.Method.(*jwt.SigningMethodHMAC)
						if !ok {
							return nil, nil
						}
						return []byte("c3VwZXJzZWNyZXRzdXBlcmR1cGVyc2VjcmV0"),
							nil
					})
				if error != nil {
					http.Error(w, fmt.Sprintf("ERROR - FAILED - %s",
						error.Error()), 400)
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(models.Exception{
						Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(models.Exception{
				Message: "An authorization header is required"})
		}
	})
}

/*
GetToken - receive a token for authentication, returns 400 for errors
*/
func (c *Controller) GetToken(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var login models.Login
	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		error := models.RespError{
			Error: "Failed request, could not parse json. " +
				"Looking for username and password",
		}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	//connect to database
	db, err := c.Session.Connect()
	if err != nil {
		error := models.RespError{
			Error: "Failed to connect, cannot reach database",
		}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	defer db.Close()

	err = login.ValidateLogin(db)
	if err != nil {
		error := models.RespError{
			Error: "Error during login validation. " +
				"Please check credentials",
		}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       login.ID,
		"username": login.Username,
	})
	tokenString, err := token.SignedString(
		[]byte("c3VwZXJzZWNyZXRzdXBlcmR1cGVyc2VjcmV0"))
	if err != nil {
		error := models.RespError{
			Error: "Make sure you have permissions to use this route. " +
				"Failed  to get token",
		}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	userid := login.ID
	login = models.Login{}
	w.WriteHeader(http.StatusOK)
	c.Logger.Logging(req, 200)
	json.NewEncoder(w).Encode(models.JwtToken{
		Token: tokenString, UserProfileID: userid,
	})
}
