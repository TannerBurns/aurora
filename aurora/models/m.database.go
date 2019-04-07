package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Connection struct {
	LiteConfig *LiteConfig
}

func NewSession(lc *LiteConfig) *Connection {
	connection := &Connection{}
	connection.LiteConfig = lc
	return connection
}

// Connect to database
func (conn *Connection) Connect() (db *sql.DB, err error) {
	if conn.LiteConfig.Config["postgres"]["password"] == "" {
		db, err = sql.Open("postgres",
			fmt.Sprintf("host=%s dbname=%s user=%s sslmode=disable",
				conn.LiteConfig.Config["postgres"]["host"],
				conn.LiteConfig.Config["postgres"]["database"],
				conn.LiteConfig.Config["postgres"]["user"]))
	} else {
		db, err = sql.Open("postgres",
			fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
				conn.LiteConfig.Config["postgres"]["host"],
				conn.LiteConfig.Config["postgres"]["database"],
				conn.LiteConfig.Config["postgres"]["user"],
				conn.LiteConfig.Config["postgres"]["password"]))
	}
	return
}
