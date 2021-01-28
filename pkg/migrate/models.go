package migrate

import (
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	Id              int
	Version         string
	MigrationStatus string
	LastUpdated     time.Time
}

type DatabaseConnection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type MigrationStore interface {
	Init() error
	GetCurrentSchemaVersion() (string, error)
	MigrateSchema(version, sql string) error
}

type SchemaMigrationStore struct {
	Db        DatabaseConnection
	TableName string
}

func (s *SchemaMigrationStore) Init() error {
	query := `CREATE TABLE IF NOT EXISTS %s(
		id SERIAL PRIMARY KEY,
		version VARCHAR(16) NOT NULL,
		migration_status VARCHAR(16) DEFAULT 'in progress',
		last_updated TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	_, err := s.Db.Exec(fmt.Sprintf(query, s.TableName))
	if err != nil {
		return err
	}

	query = `INSERT INTO %s(version, migration_status) VALUES('000', 'success')`
	_, err = s.Db.Exec(fmt.Sprintf(query, s.TableName))
	if err != nil {
		return err
	}

	return nil
}

func (s *SchemaMigrationStore) initialized() bool {
	query := fmt.Sprintf("SELECT COUNT(*) from %s", s.TableName)

	_, err := s.Db.Exec(query)
	if err != nil {
		return false
	}

	return true
}

func (s *SchemaMigrationStore) GetCurrentSchemaVersion() (string, error) {
	if !s.initialized() {
		return "", ErrDatabaseNotInitialized
	}

	query := "SELECT version FROM %s WHERE id=(SELECT MAX(id) FROM %s)"
	result := s.Db.QueryRow(fmt.Sprintf(query, s.TableName, s.TableName))

	var currentVersion string
	err := result.Scan(&currentVersion)
	if err != nil {
		return "", err
	}

	return currentVersion, nil
}

func (s *SchemaMigrationStore) startMigration(version string) error {
	query := fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", s.TableName)
	_, err := s.Db.Exec(query, version)
	if err != nil {
		return err
	}

	return nil
}

func (s *SchemaMigrationStore) endMigration(version string, migrationSuccessful bool) error {
	var migrationStatus string
	if migrationSuccessful {
		migrationStatus = "success"
	} else {
		migrationStatus = "failure"
	}

	query := fmt.Sprintf("UPDATE %s SET migration_status=$1, last_updated=NOW() WHERE version=$2", s.TableName)
	_, err := s.Db.Exec(query, migrationStatus, version)
	if err != nil {
		return err
	}

	return nil
}

func (s *SchemaMigrationStore) MigrateSchema(version, sql string) error {
	err := s.startMigration(version)
	if err != nil {
		return err
	}

	_, migrationErr := s.Db.Exec(sql)
	var migrationSuccessful bool
	if migrationErr != nil {
		migrationSuccessful = false
	} else {
		migrationSuccessful = true
	}

	err = s.endMigration(version, migrationSuccessful)
	if err != nil {
		return err
	}

	return nil
}

func NewSchemaMigrationStore(db DatabaseConnection) *SchemaMigrationStore {
	sm := SchemaMigrationStore{
		Db:        db,
		TableName: "pgm_schema_migration",
	}

	return &sm
}
