package models

import (
	"database/sql"
	"errors"
)

type Login struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (login *Login) UpdatePassword(db *sql.DB, oldpassword string,
	newpassword string, newpasswordcheck string) (err error) {
	query := `
		UPDATE userlogin
		SET password=crypt($1, gen_salt('bf'))
		WHERE username=$2 AND password=crypt(oldpassword, gen_salt('bf'));
	`

	if newpassword == newpasswordcheck {
		_, err = db.Query(query, newpassword, login.Username)
		if err != nil {
			return
		}
	} else {
		err = errors.New("New passwords do not match")
		return
	}
	return
}

func (login *Login) ValidateLogin(db *sql.DB) (err error) {
	query := `
		SELECT id FROM userlogin 
		WHERE username=$1 AND password=crypt($2, password);
	`
	err = db.QueryRow(query, login.Username, login.Password).Scan(
		&login.ID)
	return
}
