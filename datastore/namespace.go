package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
)

type Namespace string

func getNamespaces(ctx context.Context) ([]*Namespace, error) {

	id, err := setEnvironment()
	if err != nil {
		return nil, xerrors.Errorf("setEnvironment() error: %w", err)
	}

	//namespace 指定の場合
	cli, err := datastore.NewClient(ctx, id)
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
