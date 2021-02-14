package datastore

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

type Kind struct {
	Name       string
	Properties []*Property
}

func (k Kind) String() string {
	var b strings.Builder
	b.WriteString(k.Name)
	for _, prop := range k.Properties {
		b.WriteString("\n  " + prop.String())
	}
	return b.String()

}

type Property struct {
	Name string
	Type []string
}

func (p Property) String() string {
	return fmt.Sprintf("%s: %v", p.Name, p.Type)
}

type PropertyValue struct {
	Repr []string `datastore:"property_representation"`
}

func getKinds(ctx context.Context, s *config.Project) ([]*Kind, error) {

	//namespace 指定の場合
	cli, err := datastore.NewClient(ctx, s.ProjectID)
	if err != nil {
		return nil, xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	q := datastore.NewQuery("__property__")

	var props []PropertyValue
	keys, err := cli.GetAll(ctx, q, &props)
	if err != nil {
		return nil, xerrors.Errorf("property GetAll() error: %w", err)
	}

	kindMap := make(map[string][]*Property)

	for idx, key := range keys {
		prop := props[idx]
		p := Property{key.Name, prop.Repr}
		name := key.Parent.Name
		kindMap[name] = append(kindMap[name], &p)
	}

	var kinds []*Kind
	for key, elm := range kindMap {
		kind := Kind{key, elm}
		kinds = append(kinds, &kind)
	}

	return kinds, nil
}

func ViewAllKind(ctx context.Context, s *config.Project) error {
	kinds, err := getKinds(context.Background(), s)
	if err != nil {
		return xerrors.Errorf("datastore.getKinds() error: %w", err)
	}

	for _, kind := range kinds {
		fmt.Println("=================")
		fmt.Println(kind)
	}
	return nil
}

func RemoveAllKind(ctx context.Context, s *config.Project) error {

	kinds, err := getKinds(ctx, s)
	if err != nil {
		return xerrors.Errorf("getKinds() error: %w", err)
	}

	cli, err := datastore.NewClient(ctx, s.ProjectID)
	if err != nil {
		return xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	_, err = cli.RunInTransaction(ctx, func(tx *datastore.Transaction) error {

		for _, kind := range kinds {

			name := kind.Name

			q := datastore.NewQuery(name).KeysOnly()
			keys, err := cli.GetAll(ctx, q, nil)
			if err != nil {
				return xerrors.Errorf("GetAll()[%s]: %w", name, err)
			}

			err = tx.DeleteMulti(keys)
			if err != nil {
				return xerrors.Errorf("delete multi[%s]: %w", name, err)
			}
		}

		return nil
	})

	if err != nil {
		return xerrors.Errorf("remove kind error: %w", err)
	}
	return nil
}
