package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type NullInt64 struct {
	sql.NullInt64
}
type NullString struct {
	sql.NullString
}
type NullTime struct {
	mysql.NullTime
}

func (ns *NullString) MarshalJSON() (resp []byte, err error) {
	if !ns.Valid {
		resp = []byte("null")
	}
	resp, err = json.Marshal(ns.String)
	return
}

func (ns *NullString) UnmarshalJSON(b []byte) (err error) {
	err = json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return
}

func (ni *NullInt64) MarshalJSON() (resp []byte, err error) {
	if !ni.Valid {
		resp = []byte("null")
	}
	resp, err = json.Marshal(ni.Int64)
	return
}

func (ni *NullInt64) UnmarshalJSON(b []byte) (err error) {
	err = json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func (nt *NullTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		nt.Valid = false
		return
	}

	nt.Time = x
	nt.Valid = true
	return
}
