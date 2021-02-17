package datastore

import (
	"os"

	"github.com/secondarykey/dem/config"
)

const (
	DatastoreEmulatorHostEnv = "DATASTORE_EMULATOR_HOST"
	DatastoreProjectIDEnv    = "DATASTORE_PROJECT_ID"
	DatastoreDatasetEnv      = "DATASTORE_DATASET"
)

func setEnv(p *config.Project) {
	os.Setenv(DatastoreEmulatorHostEnv, p.Endpoint)
	os.Setenv(DatastoreProjectIDEnv, p.ProjectID)
	os.Setenv(DatastoreDatasetEnv, p.ProjectID)
	return
}
