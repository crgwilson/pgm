package migrate

import (
	"sort"

	"github.com/crgwilson/pgm/pkg/logger"
)

const schemaVersionTableName = "pgm_schema_migration"

type MigrationManager struct {
	Datastore        MigrationStore
	SchemaVersions   []string
	SchemaVersionMap map[string]*SchemaVersion
	Logger           logger.CliLogger
}

func (m *MigrationManager) InitDb() error {
	m.Logger.Debug("Preparing to initialized " + schemaVersionTableName + " table in target database")
	err := m.Datastore.Init()
	if err != nil {
		m.Logger.Error("An error has occurred while trying to create migration table")
		return err
	}
	m.Logger.Debug("Migration table created successfully")

	return nil
}

func (m *MigrationManager) CurrentVersion() (string, error) {
	currentVersion, err := m.Datastore.GetCurrentSchemaVersion()
	if err != nil {
		return "", err
	}

	return currentVersion, nil
}

func (m *MigrationManager) LowestAvailableVersion() string {
	version := m.SchemaVersions[0]
	return version
}

func (m *MigrationManager) HighestAvailableVersion() string {
	version := m.SchemaVersions[len(m.SchemaVersions)-1]
	return version
}

func (m *MigrationManager) isKnownVersion(version string) bool {
	_, ok := m.SchemaVersionMap[version]
	return ok
}

func (m *MigrationManager) addSchemaVersion(schema *SchemaVersion) error {
	if m.isKnownVersion(schema.Version) {
		return ErrSchemaVersionAlreadyDefined
	}

	// Add schema pointed to map for easy access
	m.SchemaVersionMap[schema.Version] = schema

	// Add schema version name to slice to maintain proper ordering
	newVersionSlice := append(m.SchemaVersions, schema.Version)

	sort.Strings(newVersionSlice)

	m.SchemaVersions = newVersionSlice

	return nil
}

func (m *MigrationManager) getVersionIndex(version string) (int, error) {
	for i, v := range m.SchemaVersions {
		if v == version {
			return i, nil
		}
	}

	return 0, ErrSchemaVersionUnknown
}

func (m *MigrationManager) getNextStepUp() (SchemaVersion, error) {
	version, err := m.CurrentVersion()
	if err != nil {
		return SchemaVersion{}, err
	}

	var nextSchemaVersion *SchemaVersion
	if version == "000" {
		nextVersion := m.SchemaVersions[0]
		nextSchemaVersion = m.SchemaVersionMap[nextVersion]
	} else {
		versionIndex, err := m.getVersionIndex(version)
		if err != nil {
			return SchemaVersion{}, err
		}
		lastIndex := len(m.SchemaVersions) - 1
		if versionIndex == lastIndex {
			return SchemaVersion{}, ErrNoNextStep
		}

		nextVersion := m.SchemaVersions[versionIndex+1]
		nextSchemaVersion = m.SchemaVersionMap[nextVersion]
	}

	return *nextSchemaVersion, nil
}

func (m *MigrationManager) getNextStepDown() (SchemaVersion, error) {
	version, err := m.CurrentVersion()
	if err != nil {
		return SchemaVersion{}, err
	}

	versionIndex, err := m.getVersionIndex(version)
	if err != nil {
		return SchemaVersion{}, err
	}

	if versionIndex == 0 {
		return SchemaVersion{}, ErrNoNextStep
	}

	nextVersion := m.SchemaVersions[versionIndex-1]
	nextVersionSchema := m.SchemaVersionMap[nextVersion]

	return *nextVersionSchema, nil
}

func (m *MigrationManager) Up(targetVersion string) error {
	version, err := m.CurrentVersion()
	if err != nil {
		return err
	}

	if version == targetVersion {
		m.Logger.Info("Reached target version " + targetVersion)
		return nil
	}

	next, err := m.getNextStepUp()
	if err != nil {
		return err
	}

	m.Logger.Info("Beginning schema migration from version " + version + " to " + next.Version)
	err = m.Datastore.MigrateSchema(next.Version, next.Up)
	if err != nil {
		return err
	}

	// We want to keep going step by step until we reach our target version
	err = m.Up(targetVersion)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (m *MigrationManager) Down(targetVersion string) error {
	version, err := m.CurrentVersion()
	if err != nil {
		return err
	}

	if version == targetVersion {
		m.Logger.Info("Reached target version " + targetVersion)
		return nil
	}

	down := m.SchemaVersionMap[version]
	next, err := m.getNextStepDown()
	if err != nil {
		return err
	}

	m.Logger.Info("Beginning schema migration from version " + version + " to " + next.Version)
	err = m.Datastore.MigrateSchema(next.Version, down.Down)
	if err != nil {
		return err
	}

	err = m.Down(targetVersion)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (m *MigrationManager) RegisterMigrationPath(migrationPath MigrationPath) error {
	schema, versionExists := m.SchemaVersionMap[migrationPath.Version]
	if !versionExists {
		schema = NewSchemaVersion(migrationPath.Version)
		m.addSchemaVersion(schema)
		m.Logger.Debug("Schema for version " + schema.Version + " does not already exist, creating a new schema definition")
	}

	err := schema.SetAction(migrationPath.Action, migrationPath.Sql())
	m.Logger.Debug("Registered action " + migrationPath.Action + " for schema version " + migrationPath.Version)
	if err != nil {
		return err
	}

	return nil
}

func NewMigrationManager(db MigrationStore, l logger.CliLogger) *MigrationManager {
	migrator := MigrationManager{
		Datastore:        db,
		SchemaVersions:   make([]string, 0),
		SchemaVersionMap: make(map[string]*SchemaVersion),
		Logger:           l,
	}

	return &migrator
}
