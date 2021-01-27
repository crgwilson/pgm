package migrate

import "errors"

var ErrInvalidAction = errors.New("Schema migration 'action' must be set to either 'up' or 'down'")

type SchemaVersion struct {
	Version string
	Up      string
	Down    string
}

func (s *SchemaVersion) SetAction(action, sqlText string) error {
	switch action {
	case "up":
		s.Up = sqlText
	case "down":
		s.Down = sqlText
	default:
		return ErrInvalidAction
	}
	return nil
}

func NewSchemaVersion(version string) *SchemaVersion {
	sv := SchemaVersion{
		Version: version,
	}
	return &sv
}
