package migrate

import (
	"testing"
)

const testVersionNumber = "001"
const testUpSql = "CREATE TABLE IF NOT EXISTS test(id SERIAL PRIMARY KEY)"
const testDownSql = "DROP TABLE test"

func TestSchemaVersion(t *testing.T) {
	cases := []struct {
		Name            string
		Input           MigrationPath
		ExpectedVersion string
		ExpectedUpSql   string
		ExpectedDownSql string
		ExpectedError   error
	}{
		{
			"creating new schema version without adding actions",
			MigrationPath{
				Version: "test1",
			},
			"test1",
			"",
			"",
			nil,
		},
		{
			"creating new schema version and adding 'up' action",
			MigrationPath{
				Version: "test2",
				Action:  "up",
				Raw:     []byte(testUpSql),
			},
			"test2",
			testUpSql,
			"",
			nil,
		},
		{
			"creating new schema version and adding 'down' action",
			MigrationPath{
				Version: "test3",
				Action:  "down",
				Raw:     []byte(testDownSql),
			},
			"test3",
			"",
			testDownSql,
			nil,
		},
		{
			"creating new schema version and adding 'invalid' action",
			MigrationPath{
				Version: "test4",
				Action:  "invalid",
			},
			"test4",
			"",
			"",
			ErrInvalidAction,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			testSchemaVersion := NewSchemaVersion(test.Input.Version)
			if testSchemaVersion.Version != test.ExpectedVersion {
				t.Errorf("got %q, want %q", testSchemaVersion.Version, test.ExpectedVersion)
			}

			if test.Input.Action != "" {
				err := testSchemaVersion.SetAction(test.Input.Action, test.Input.Sql())
				if err != test.ExpectedError {
					t.Errorf("got %v, want %v", err, test.ExpectedError)
				}
			}

			if testSchemaVersion.Up != test.ExpectedUpSql {
				t.Errorf("got %q, want %q", testSchemaVersion.Up, test.ExpectedUpSql)
			}

			if testSchemaVersion.Down != test.ExpectedDownSql {
				t.Errorf("got %q, want %q", testSchemaVersion.Down, test.ExpectedDownSql)
			}
		})
	}
}
