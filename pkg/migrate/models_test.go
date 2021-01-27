package migrate

type MockMigrationStore struct {
	currentVersion *Migration
	migrations     []Migration
}

func (m *MockMigrationStore) Init() error {
	return nil
}

func (m *MockMigrationStore) GetCurrentSchemaVersion() (string, error) {
	return m.currentVersion.Version, nil
}

func (m *MockMigrationStore) MigrateSchema(version, sql string) error {
	newMigration := Migration{
		Version: version,
	}
	m.migrations = append(m.migrations, newMigration)
	m.currentVersion = &newMigration

	return nil
}

func NewMockMigrationStore() *MockMigrationStore {
	migration := Migration{
		Version: "000",
	}

	m := MockMigrationStore{
		currentVersion: &migration,
		migrations:     []Migration{migration},
	}

	return &m
}
