package datastore

import (
	"os"

	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

const (
	DatastoreEmulatorHostEnv = "DATASTORE_EMULATOR_HOST"
	DatastoreProjectIDEnv    = "DATASTORE_PROJECT_ID"
	DatastoreDatasetEnv      = "DATASTORE_DATASET"
)

func setEnvironment() (string, error) {
	p := config.GetCurrentProject()
	if p == nil {
		return "", xerrors.Errorf("Current Project empty.")
	}
	os.Setenv(DatastoreEmulatorHostEnv, p.Endpoint)
	os.Setenv(DatastoreProjectIDEnv, p.ProjectID)
	os.Setenv(DatastoreDatasetEnv, p.ProjectID)
	return p.ProjectID, nil
}
