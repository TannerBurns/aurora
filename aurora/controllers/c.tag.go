package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"../models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

/*
CreateNewTag - create a new tag for comment.
	returns 400, 401, 404 for errors
*/
func (c *Controller) CreateTag(w http.ResponseWriter, req *http.Request) {
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

	//get task id, should be in route
	params := mux.Vars(req)
	tid, err := strconv.Atoi(params["tid"])
	if err != nil {
		error := models.RespError{Error: "Task id is required in route" +
			" in route"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	//get comment id, should be in route
	cid, err := strconv.Atoi(params["cid"])
	if err != nil {
		error := models.RespError{Error: "comment id is required in route" +
			" in route"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}

	//get body, should be in the format of a tag
	var tag *models.Tag
	err = json.NewDecoder(req.Body).Decode(&tag)
	if err != nil {
		error := models.RespError{
			Error: "Failed to parse request." +
				" Please make sure request is valid format"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	tag.OwnerID = userprofile.ID
	tag.CommentID = cid

	//create new tag
	err = userprofile.CreateTag(db, tag)
	if err != nil {
		error := models.RespError{Error: "Failed to create a new tag"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//get task to return
	t, err := userprofile.GetTask(db, tid)
	if err != nil {
		error := models.RespError{Error: "Failed to get task"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//send response, 200
	w.WriteHeader(http.StatusOK)
	c.Logger.Logging(req, 200)
	json.NewEncoder(w).Encode(t)
	return
}

/*
DeleteTag - delete a tag for comment
	returns 400, 401, 404 for errors
*/
func (c *Controller) DeleteTag(w http.ResponseWriter, req *http.Request) {
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

	//get task id, should be in route
	params := mux.Vars(req)
	tid, err := strconv.Atoi(params["tid"])
	if err != nil {
		error := models.RespError{Error: "Task id is required in route" +
			" in route"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}
	//get tag id, should be in route
	taid, err := strconv.Atoi(params["taid"])
	if err != nil {
		error := models.RespError{Error: "Tag id is required in route" +
			" in route"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 400)
		c.Logger.Logging(req, 400)
		return
	}

	//delete tag
	err = userprofile.DeleteTag(db, taid)
	if err != nil {
		error := models.RespError{Error: "Failed to delete tag"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//get task to return
	t, err := userprofile.GetTask(db, tid)
	if err != nil {
		error := models.RespError{Error: "Failed to get task"}
		resp, _ := json.Marshal(error)
		http.Error(w, string(resp), 404)
		c.Logger.Logging(req, 404)
		return
	}
	//send response, 200
	w.WriteHeader(http.StatusOK)
	c.Logger.Logging(req, 200)
	json.NewEncoder(w).Encode(t)
	return
}
