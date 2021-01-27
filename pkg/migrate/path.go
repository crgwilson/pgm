package migrate

import (
	"errors"
	"strings"
)

var ErrInvalidFile = errors.New("Provided file name not formatted as expected")

type MigrationPath struct {
	Version string
	Action  string
	Raw     []byte
}

func (p MigrationPath) Sql() string {
	s := string(p.Raw)
	return s
}

func ParseSqlFile(sqlFileName string, sqlFileContents []byte) (MigrationPath, error) {
	splitFileName := strings.Split(sqlFileName, ".")

	if len(splitFileName) != 3 {
		return MigrationPath{}, ErrInvalidFile
	}

	if splitFileName[2] != "sql" {
		return MigrationPath{}, ErrInvalidFile
	}

	if splitFileName[1] != "up" && splitFileName[1] != "down" {
		return MigrationPath{}, ErrInvalidFile
	}

	parsed := MigrationPath{
		Version: splitFileName[0],
		Action:  splitFileName[1],
		Raw:     sqlFileContents,
	}

	return parsed, nil
}
