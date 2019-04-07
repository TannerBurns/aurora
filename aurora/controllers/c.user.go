package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"../models"
	"github.com/gorilla/context"
)

/*
CreateUser - create a new user. returns 400, 401, 404 for errors
*/
func (c *Controller) CreateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//connect to database
	db, err := c.Session.Connect()
	if err != nil {
		error := models.RespError{
			Error: "Failed to connect, cannot reach database"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	defer db.Close()

	//create user object to hold user information in body
	var userprofile models.UserProfile
	err = json.NewDecoder(req.Body).Decode(&userprofile)
	if err != nil {
		error := models.RespError{
			Error: "Failed to parse request." +
				" Please make sure request is valid format"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//create a new user
	err = userprofile.CreateUser(db)
	if err != nil {
		fmt.Println(err)
		error := models.RespError{Error: "Failed to create a new user"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//send response, 200
	w.WriteHeader(http.StatusOK)
	c.Logger.Logging(req, 200)
	json.NewEncoder(w).Encode(userprofile)
	return
}

/*
ReadUser - get user information. returns 400, 401, 404 for errors
*/
func (c *Controller) ReadUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//connect to database
	db, err := c.Session.Connect()
	if err != nil {
		error := models.RespError{
			Error: "Failed to connect, cannot reach database"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	defer db.Close()

	//get user ID from decoded jwt
	dec := context.Get(req, "decoded").(jwt.MapClaims)
	userprofile := &models.UserProfile{ID: int(dec["id"].(float64))}

	//update user structure
	err = userprofile.ReadUser(db)
	if err != nil {
		error := models.RespError{
			Error: "Failed to get user with ID: " +
				fmt.Sprintf("%d", int(dec["id"].(float64))),
		}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//send response, 200
	w.WriteHeader(http.StatusOK)
	c.Logger.Logging(req, 200)
	json.NewEncoder(w).Encode(userprofile)
	return
}

/*
UpdateUser - update a user by id. returns 400, 401, 404 for errors
*/
func (c *Controller) UpdateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//connect to database
	db, err := c.Session.Connect()
	if err != nil {
		error := models.RespError{
			Error: "Failed to connect, cannot reach database"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	defer db.Close()

	//pull user ID from decoded jwt
	dec := context.Get(req, "decoded").(jwt.MapClaims)
	userprofile := models.UserProfile{ID: int(dec["id"].(float64))}
	err = json.NewDecoder(req.Body).Decode(&userprofile)
	if err != nil {
		error := models.RespError{
			Error: "Failed to parse request." +
				" Please make sure request is valid format"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//update user in database and current structure
	err = userprofile.UpdateUser(db)
	if err != nil {
		error := models.RespError{
			Error: "Failed to get user with ID: " +
				fmt.Sprintf("%d", int(dec["id"].(float64))),
		}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//send response, 200
	w.WriteHeader(http.StatusOK)
	c.Logger.Logging(req, 200)
	json.NewEncoder(w).Encode(userprofile)
	return
}
