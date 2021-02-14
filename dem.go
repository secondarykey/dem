package dem

import (
	"context"
	"fmt"
	"os"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"
	"golang.org/x/xerrors"
)

const (
	DefaultNamespace = "[default]"
	DefaultProjectID = "[empty]"
)

const (
	DatastoreEmulatorHostEnv = "DATASTORE_EMULATOR_HOST"
	DatastoreProjectIDEnv    = "DATASTORE_PROJECT_ID"
	DatastoreDatasetEnv      = "DATASTORE_DATASET"
	DefaultEndpoint          = "localhost:8081"
)

func setEnv() *config.Project {
	var s config.Project
	s.ProjectID = "section"
	s.Endpoint = DefaultEndpoint
	os.Setenv(DatastoreEmulatorHostEnv, s.Endpoint)
	os.Setenv(DatastoreProjectIDEnv, s.ProjectID)
	os.Setenv(DatastoreDatasetEnv, s.ProjectID)
	return &s
}

func Remove(names ...string) error {

	s := setEnv()
	ctx := context.Background()
	err := datastore.RemoveAllKind(ctx, s)
	if err != nil {
		return xerrors.Errorf("datastore.RemoveAllKind() error: %w", err)
	}
	return nil
}

func ViewSchema() error {
	s := setEnv()
	err := datastore.ViewAllKind(context.Background(), s)
	if err != nil {
		return xerrors.Errorf("datastore.ViewAllKind() error: %w", err)
	}
	return nil
}

func ViewEntity(kind string) error {

	s := setEnv()

	entities, err := datastore.GetEntities(context.Background(), s, kind)
	if err != nil {
		return xerrors.Errorf("GetEntities() error: %w", err)
	}

	for _, elm := range entities {
		fmt.Println(elm)
	}

	return nil
}
