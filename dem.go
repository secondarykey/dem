package dem

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

func Listen() error {

	s := mux.NewRouter()
	s.HandleFunc("/", indexHandler)
	s.HandleFunc("/{id}/", kindHandler)
	s.HandleFunc("/{id}/{kind}", entityHandler)

	return http.ListenAndServe(":8088", s)
}

func setEnv() *config.Project {
	var s config.Project
	s.ProjectID = "section"
	s.Endpoint = DefaultEndpoint
	os.Setenv(DatastoreEmulatorHostEnv, s.Endpoint)
	os.Setenv(DatastoreProjectIDEnv, s.ProjectID)
	os.Setenv(DatastoreDatasetEnv, s.ProjectID)
	return &s
}

func getKinds(names ...string) ([]*datastore.Kind, error) {
	s := setEnv()
	ctx := context.Background()

	kinds, err := datastore.GetKinds(ctx, s, names...)
	if err != nil {
		return nil, xerrors.Errorf("datastore.GetKinds() error: %w", err)
	}
	return kinds, nil
}

func RemoveEntity(names ...string) error {

	kinds, err := getKinds(names...)
	if err != nil {
		return xerrors.Errorf("getKinds() error: %w", err)
	}

	s := setEnv()
	ctx := context.Background()
	for _, kind := range kinds {
		err := datastore.RemoveKind(ctx, s, kind.Name)
		if err != nil {
			return xerrors.Errorf("datastore.RemoveAllKind() error: %w", err)
		}
	}
	return nil
}

func ViewKind(names ...string) error {
	kinds, err := getKinds(names...)
	if err != nil {
		return xerrors.Errorf("getKinds() error: %w", err)
	}

	for _, kind := range kinds {
		fmt.Println(kind)
	}

	return nil
}

func ViewEntity(names ...string) error {

	s := setEnv()

	kinds, err := getKinds(names...)
	if err != nil {
		return xerrors.Errorf("getKinds() error: %w", err)
	}

	for _, kind := range kinds {
		entities, err := datastore.GetEntities(context.Background(), s, kind.Name)
		if err != nil {
			return xerrors.Errorf("GetEntities() error: %w", err)
		}

		fmt.Println("################### " + kind.Name)
		for _, elm := range entities {
			fmt.Println(elm)
		}
	}

	return nil
}
