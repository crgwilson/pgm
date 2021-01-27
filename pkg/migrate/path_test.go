package migrate

import (
	"testing"
)

const testNameUp = "001.up.sql"
const testSqlUp = `CREATE TABLE IF NOT EXISTS test_table(
	id SERIAL PRIMARY KEY,
	name VARCHAR(10) NOT NULL
);
`

const testNameDown = "002.down.sql"
const testSqlDown = `DROP TABLE some_other_table;`

const testNameNoExtension = "somestring"
const testNameNotSql = "somestring.up.txt"
const testNameInvalidAction = "somestring.else.sql"

type TestParseSqlFileInput struct {
	FileName     string
	FileContents []byte
}

func TestParseSqlFile(t *testing.T) {
	cases := []struct {
		Name            string
		Input           TestParseSqlFileInput
		ExpectedVersion string
		ExpectedAction  string
		ExpectedSql     string
		ExpectedError   error
	}{
		{
			"version 001, action up",
			TestParseSqlFileInput{
				FileName:     testNameUp,
				FileContents: []byte(testSqlUp),
			},
			"001",
			"up",
			testSqlUp,
			nil,
		},
		{
			"version 002, action down",
			TestParseSqlFileInput{
				FileName:     testNameDown,
				FileContents: []byte(testSqlDown),
			},
			"002",
			"down",
			testSqlDown,
			nil,
		},
		{
			"file name with no file extension",
			TestParseSqlFileInput{
				FileName:     testNameNoExtension,
				FileContents: []byte("this does not matter"),
			},
			"",
			"",
			"",
			ErrInvalidFile,
		},
		{
			"file name with non-sql extension",
			TestParseSqlFileInput{
				FileName:     testNameNotSql,
				FileContents: []byte("this does not matter"),
			},
			"",
			"",
			"",
			ErrInvalidFile,
		},
		{
			"file name with invalid action",
			TestParseSqlFileInput{
				FileName:     testNameInvalidAction,
				FileContents: []byte("this does not matter"),
			},
			"",
			"",
			"",
			ErrInvalidFile,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			parsed, err := ParseSqlFile(test.Input.FileName, test.Input.FileContents)

			if err != test.ExpectedError {
				t.Errorf("got %v, want %v", err, test.ExpectedError)
			}

			if parsed.Version != test.ExpectedVersion {
				t.Errorf("got %q, want %q", parsed.Version, test.ExpectedVersion)
			}

			if parsed.Action != test.ExpectedAction {
				t.Errorf("got %q, want %q", parsed.Action, test.ExpectedAction)
			}

			got := parsed.Sql()
			if got != test.ExpectedSql {
				t.Errorf("got %q, want %q", got, test.ExpectedSql)
			}
		})
	}
}
