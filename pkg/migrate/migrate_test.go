package migrate

import (
	"testing"

	"github.com/crgwilson/pgm/pkg/logger"
	"github.com/crgwilson/pgm/pkg/mocks"
)

func TestMigrations(t *testing.T) {
	lgr := logger.CliLogger{
		Logger:   mocks.NewSpyLogger(),
		LogLevel: logger.DebugLogLevel(),
	}

	db := NewMockMigrationStore()

	testMigrator := NewMigrationManager(db, lgr)
	err := testMigrator.InitDb()
	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	currentVersion, err := testMigrator.CurrentVersion()
	if currentVersion != "000" {
		t.Errorf("got %q, want 000", currentVersion)
	}

	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	schemaOneUp := MigrationPath{
		Version: "001",
		Action:  "up",
		Raw:     []byte("001up"),
	}

	known := testMigrator.isKnownVersion("001")
	if known {
		t.Errorf("somehow found unknown version")
	}

	err = testMigrator.RegisterMigrationPath(schemaOneUp)
	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	idx, err := testMigrator.getVersionIndex("001")
	if idx != 0 {
		t.Errorf("got %d, want %d", idx, 0)
	}

	known = testMigrator.isKnownVersion("001")
	if !known {
		t.Errorf("cannot find expected migration version")
	}

	schemaTwoUp := MigrationPath{
		Version: "002",
		Action:  "up",
		Raw:     []byte("002up"),
	}

	schemaThreeUp := MigrationPath{
		Version: "003",
		Action:  "up",
		Raw:     []byte("003up"),
	}

	err = testMigrator.RegisterMigrationPath(schemaTwoUp)
	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	err = testMigrator.RegisterMigrationPath(schemaThreeUp)
	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	highest := testMigrator.HighestAvailableVersion()
	if highest != "003" {
		t.Errorf("got %q, want %q", highest, "003")
	}

	lowest := testMigrator.LowestAvailableVersion()
	if lowest != "001" {
		t.Errorf("got %q, want %q", lowest, "001")
	}

	next, err := testMigrator.getNextStepUp()
	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	if next.Version != "001" {
		t.Errorf("got %q, want %q", next.Version, "001")
	}

	err = testMigrator.Up("003")
	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	currentVersion, err = testMigrator.CurrentVersion()
	if err != nil {
		t.Errorf("got %v, want no error", err)
	}

	if currentVersion != "003" {
		t.Errorf("got %q, want %q", currentVersion, "003")
	}
}
