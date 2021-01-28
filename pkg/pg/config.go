package pg

import (
	"errors"
	"fmt"
)

var ErrMissingAddress = errors.New("Database address is required but not set")
var ErrMissingPort = errors.New("Database port number is required but not set")

type PostgresConfig struct {
	Address  string
	Port     int
	User     string
	Password string
	Database string
	SslMode  string
}

func (c PostgresConfig) ConnectionString() (string, error) {
	connString := "postgres://"

	if c.User != "" {
		connString = connString + c.User
	}

	if c.User != "" && c.Password != "" {
		connString = connString + ":" + c.Password + "@"
	}

	if c.Address == "" {
		return "", ErrMissingAddress
	}

	connString = connString + c.Address

	if c.Port == 0 {
		return "", ErrMissingPort
	}

	connString = fmt.Sprintf("%s:%d/", connString, c.Port)

	if c.Database != "" {
		connString = connString + c.Database
	}

	if c.SslMode != "" {
		connString = connString + "?sslmode=" + c.SslMode
	}

	return connString, nil
}

