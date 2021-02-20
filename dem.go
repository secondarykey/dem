package dem

import (
	"context"
	"fmt"
	"net/http"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"
	"github.com/secondarykey/dem/handler"

	"golang.org/x/xerrors"
)

const (
	DefaultNamespace = "[default]"
	DefaultProjectID = "[empty]"
)

func SetConsoleOptions(opts ...config.ConsoleOption) error {
	err := config.SetConsole(opts)
	if err != nil {
		return xerrors.Errorf("config.SetConsole() error: %w", err)
	}
	return nil
}

func Listen(opts ...config.ViewerOption) error {

	err := config.SetViewer(opts)
	if err != nil {
		return xerrors.Errorf("config.SetViewer() error: %w", err)
	}

	err = config.LoadSetting()
	if err != nil {
		return xerrors.Errorf("config.LoadSetting() error: %w", err)
	}

	conf := config.GetViewer()
	err = handler.Register()
	if err != nil {
		return xerrors.Errorf("handler.Register() error: %w", err)
	}

	server := fmt.Sprintf(":%d", conf.Port)
	fmt.Printf("Listen HTTP Server[%s]\n", server)

	return http.ListenAndServe(server, nil)
}

func getKinds(names ...string) ([]*datastore.Kind, error) {

	ctx := context.Background()
	kinds, err := datastore.GetKinds(ctx, names...)
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

	ctx := context.Background()
	for _, kind := range kinds {
		err := datastore.RemoveKind(ctx, kind.Name)
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
		fmt.Println("======================================")
		fmt.Println(kind)
	}

	return nil
}

func ViewEntity(names ...string) error {

	kinds, err := getKinds(names...)
	if err != nil {
		return xerrors.Errorf("getKinds() error: %w", err)
	}

	for _, kind := range kinds {
		entities, err := datastore.GetEntities(context.Background(), kind.Name)
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
