package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

/*
CreateNewComment - create a new comment for task.
	returns 400, 401, 404 for errors
*/
func (c *Controller) CreateNewComment(w http.ResponseWriter,
	req *http.Request) {
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

	//get id, should be in route
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		error := models.RespError{Error: "Task id is required in route"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}

	//get body, should be in the format of a comment
	var comment *models.Comment
	err = json.NewDecoder(req.Body).Decode(&comment)
	if err != nil {
		error := models.RespError{
			Error: "Failed to parse request." +
				" Please make sure request is valid format"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	comment.OwnerID = userprofile.ID
	comment.TaskID = id

	//create new comment
	err = userprofile.CreateNewComment(db, comment)
	if err != nil {
		fmt.Println(err)
		error := models.RespError{Error: "Failed to create a new comment"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}

	//send response, 200
	w.WriteHeader(http.StatusOK)
	c.Logger.Logging(req, 200)
	json.NewEncoder(w).Encode(comment)
	return
}
