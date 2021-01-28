package pg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const postgresDriverName = "postgres"

func OpenDb(connConfig PostgresConfig) (*sql.DB, error) {
	connString, err := connConfig.ConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(postgresDriverName, connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
