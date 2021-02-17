package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

type Namespace string

func getNamespaces(ctx context.Context, p *config.Project) ([]*Namespace, error) {

	setEnv(p)

	//namespace 指定の場合
	cli, err := datastore.NewClient(ctx, p.ProjectID)
	if err != nil {
		return nil, xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	q := datastore.NewQuery("__namespace__").KeysOnly()

	keys, err := cli.GetAll(ctx, q, nil)
	if err != nil {
		return nil, xerrors.Errorf("Namespace GetAll() error: %w", err)
	}

	var rtn []*Namespace
	for _, key := range keys {
		space := Namespace(key.Name)
		rtn = append(rtn, &space)
	}

	return rtn, nil
}
