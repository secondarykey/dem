package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

func GetNamespaces(ctx context.Context) ([]string, error) {

	id, err := setEnvironment()
	if err != nil {
		return nil, xerrors.Errorf("setEnvironment() error: %w", err)
	}

	cli, err := datastore.NewClient(ctx, id)
	if err != nil {
		return nil, xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	q := datastore.NewQuery("__namespace__").KeysOnly()

	keys, err := cli.GetAll(ctx, q, nil)
	if err != nil {
		return nil, xerrors.Errorf("Namespace GetAll() error: %w", err)
	}

	var rtn []string
	for _, key := range keys {
		space := key.Name
		if space == "" {
			space = config.DefaultNamespace
		}
		rtn = append(rtn, space)
	}
	return rtn, nil
}
