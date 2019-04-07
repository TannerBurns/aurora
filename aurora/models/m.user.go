package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type UserProfile struct {
	ID        int        `json:"id"`
	Created   *time.Time `json:"created"`
	First     NullString `json:"first"`
	Last      NullString `json:"last"`
	Alias     NullString `json:"alias"`
	Email     NullString `json:"email"`
	Phone     NullInt64  `json:"phone"`
	Birthday  NullTime   `json:"birthday"`
	State     NullString `json:"state"`
	Notes     NullString `json:"notes"`
	Picture   NullString `json:"picture"`
	UserLogin *UserLogin `json:"login"`
}

type UserLogin struct {
	ID            int        `json:"id"`
	Created       *time.Time `json:"created"`
	UserProfileID int        `json:"userprofile_id"`
	Username      NullString `json:"username"`
	Password      NullString
}

type UserTask struct {
	ID      int        `json:"id"`
	Created *time.Time `json:"created"`
	TaskID  int        `json:"task_id"`
	UserID  int        `json:"user_id"`
	Task    *Task      `json:"task"`
}

func (usr UserProfile) MarshalJSON() (up []byte, err error) {
	var tmp struct {
		ID       int        `json:"id"`
		Created  *time.Time `json:"created"`
		First    NullString `json:"first"`
		Last     NullString `json:"last"`
		Alias    NullString `json:"alias"`
		Email    NullString `json:"email"`
		Phone    NullInt64  `json:"phone"`
		Birthday NullTime   `json:"birthday"`
		State    NullString `json:"state"`
		Notes    NullString `json:"notes"`
		Picture  NullString `json:"picture"`
		TmpLogin struct {
			Created  *time.Time `json:"created"`
			Username NullString `json:"username"`
		} `json:"login"`
	}
	tmp.ID = usr.ID
	tmp.Created = usr.Created
	tmp.First = usr.First
	tmp.Last = usr.Last
	tmp.Alias = usr.Alias
	tmp.Email = usr.Email
	tmp.Phone = usr.Phone
	tmp.Birthday = usr.Birthday
	tmp.State = usr.State
	tmp.Notes = usr.Notes
	tmp.Picture = usr.Picture
	tmp.TmpLogin.Created = usr.UserLogin.Created
	tmp.TmpLogin.Username = usr.UserLogin.Username
	up, err = json.Marshal(&tmp)
	return
}

func (up *UserProfile) CreateUser(db *sql.DB) (err error) {
	query := `
		with upfunc as (
			INSERT INTO userprofile (first, 
									last,
									alias, 
									email, 
									phone, 
									birthday, 
									state,
									notes,
									picture)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		)
		INSERT INTO userlogin (userprofile_id, username, password)
		VALUES ((SELECT id FROM upfunc), $10, crypt($11, gen_salt('bf')))
		RETURNING id;
	`
	err = db.QueryRow(query, up.First.String, up.Last.String, up.Alias.String,
		up.Email.String, up.Phone.Int64, up.Birthday, up.State.String,
		up.Notes.String, up.Picture.String, up.UserLogin.Username.String,
		up.UserLogin.Password.String,
	).Scan(
		&up.UserLogin.UserProfileID,
	)
	if err == nil {
		up.ID = up.UserLogin.UserProfileID
	}
	return
}

/*
ReadUser - read database, table userprofile & userlogin
*/
func (up *UserProfile) ReadUser(db *sql.DB) (err error) {
	query := `
		SELECT * FROM userprofile 
		INNER JOIN userlogin ON userprofile.id=userlogin.userprofile_id 
		WHERE userprofile.id=$1;
	`
	up.UserLogin = &UserLogin{}
	err = db.QueryRow(query, up.ID).Scan(&up.ID, &up.Created, &up.First,
		&up.Last, &up.Alias, &up.Email, &up.Phone, &up.Birthday, &up.State,
		&up.Notes, &up.Picture, &up.UserLogin.ID, &up.UserLogin.Created,
		&up.UserLogin.UserProfileID, &up.UserLogin.Username,
		&up.UserLogin.Password,
	)
	return
}

/*
UpdateUser - update database, table userprofile & userlogin
*/
func (up *UserProfile) UpdateUser(db *sql.DB) (err error) {
	query := `
		UPDATE userprofile
		SET
			first=$1,
			last=$2,
			alias=$3,
			email=$4,
			phone=$5,
			birthday=$6,
			state=$7,
			notes=$8
			picture=$9
		WHERE id=$10;
		UPDATE userlogin
		SET
			username=$11
		WHERE id=$12;
	`
	up.UserLogin.ID = up.ID
	_, err = db.Query(query, up.First, up.Last, up.Alias, up.Email, up.Phone,
		up.Birthday, up.State, up.Notes, up.Picture, up.ID,
		up.UserLogin.Username, up.UserLogin.ID)
	return
}

func (up *UserProfile) CreateNewTask(db *sql.DB, t *Task) (err error) {
	query := `
		with newtask as (
			INSERT INTO task (
				owner_id,
				restricted,
				status, 
				title,
				body
			)
			VALUES($1, $2, $3, $4, $5)
			RETURNING id
		)
		INSERT INTO usertask (task_id, user_id)
		VALUES ((SELECT id FROM newtask), $6)
		RETURNING task_id;
	`
	err = db.QueryRow(query, t.OwnerID, t.Restricted, t.Status, t.Title,
		t.Body, t.OwnerID).Scan(&t.ID)
	return
}

func (up *UserProfile) GetTask(db *sql.DB, id int) (nt *Task, err error) {
	var rows *sql.Rows
	var t Task
	query := `
		SELECT * FROM task WHERE id=$1
	`
	err = db.QueryRow(query, id).Scan(&t.ID, &t.Created, &t.Modified,
		&t.OwnerID, &t.Restricted, &t.Status, &t.Title, &t.Body,
	)
	if err != nil {
		return
	}
	if t.Restricted && t.OwnerID != up.ID {
		err = errors.New("this task is not public and you are not the owner")
		return
	}
	query = `
		SELECT * FROM comment WHERE task_id=$1
	`
	rows, err = db.Query(query, id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		query := `SELECT * FROM tag WHERE comment_id=$1`
		c := Comment{}
		var tags []Tag
		err = rows.Scan(&c.ID, &c.Created, &c.Modified, &c.OwnerID, &c.TaskID,
			&c.Body)
		if err != nil {
			return
		}
		rs, err := db.Query(query, c.ID)
		if err != nil {
			return nt, err
		}
		defer rs.Close()
		for rs.Next() {
			tg := Tag{}
			err = rs.Scan(&tg.ID, &tg.Created, &tg.OwnerID, &tg.CommentID,
				&tg.Body)
			if err != nil {
				return nt, err
			}
			tags = append(tags, tg)
		}
		c.Tags = tags
		t.Comments = append(t.Comments, c)
	}
	nt = &t
	return
}

func (up *UserProfile) CreateNewComment(db *sql.DB, c *Comment) (err error) {
	query := `
		INSERT INTO comment (owner_id, task_id, body)
		VALUES ($1, $2, $3)
		RETURNING *;
	`
	err = db.QueryRow(query, c.OwnerID, c.TaskID, c.Body).Scan(&c.ID,
		&c.Created, &c.Modified, &c.OwnerID, &c.TaskID, &c.Body)

	return
}

func (up *UserProfile) CreateNewTag(db *sql.DB, t *Tag) (err error) {
	query := `
		INSERT INTO tag (owner_id, comment_id, body)
		VALUES ($1, $2, $3)
		RETURNING *;
	`
	err = db.QueryRow(query, t.OwnerID, t.CommentID, t.Body).Scan(&t.ID,
		&t.Created, &t.OwnerID, &t.CommentID, &t.Body)
	return
}

func (up *UserProfile) DeleteTag(db *sql.DB, id int) (err error) {
	query := `
		DELETE FROM tag WHERE id=$1;
	`
	_, err = db.Query(query, id)
	return
}
