package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Connection holds a connection to the database
type Connection struct {
	DB *sqlx.DB
}

// NewConnection creates a new connection to the database
func NewConnection(host, user, password, dbname, sslmode, port string) (*Connection, error) {
	var connString string

	if port == "" {
		connString = fmt.Sprintf("dbname=%s user=%s password='%s' host=%s sslmode=%s", dbname, user, password, host, sslmode)
	} else {
		connString = fmt.Sprintf("dbname=%s user=%s password='%s' host=%s sslmode=%s port=%s", dbname, user, password, host, sslmode, port)
	}

	db, error := sqlx.Connect("postgres", connString)
	return &Connection{DB: db}, err
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.DB.Close()
}
