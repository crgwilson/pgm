package migrate

import (
	"errors"
)

var ErrFailedToQuerySchemaVersion = errors.New("Unable to find current schema version in database")
var ErrDatabaseNotInitialized = errors.New("Migration table could not be found")
var ErrSchemaVersionAlreadyDefined = errors.New("Given schema version is already registered")
var ErrSchemaVersionUnknown = errors.New("Given schema has not been registered")
var ErrNoNextStep = errors.New("Given schema version has no further steps")
var ErrAlreadyReachedTargetVersion = errors.New("Requested schema version has already been deployed")
var ErrNoCurrentVersion = errors.New("No migrations have been run on this database")
