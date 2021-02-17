package datastore

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

type Kind struct {
	KeyType    KeyType
	Name       string
	Properties []*Property
}

func (k Kind) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s(Key is %s)", k.Name, k.KeyType))
	for _, prop := range k.Properties {
		b.WriteString("\n  " + prop.String())
	}
	return b.String()
}

func (k *Kind) setProperties(cli *datastore.Client, ctx context.Context) error {
	q := datastore.NewQuery(k.Name).Limit(1)
	var dst []Entity
	keys, err := cli.GetAll(ctx, q, &dst)
	if err != nil {
		return xerrors.Errorf("Property GetAll() error: %w", err)
	}

	key := keys[0]
	k.KeyType = StringKeyType
	if key.Name == "" {
		k.KeyType = IntKeyType
	}

	for _, entity := range dst {
		vals := entity.Values
		k.Properties = make([]*Property, len(vals))
		idx := 0
		for key, elm := range vals {
			v := reflect.ValueOf(elm)
			k.Properties[idx] = &Property{key, v.Type()}
			idx++
		}
	}

	sort.Slice(k.Properties, func(i, j int) bool {
		return k.Properties[i].Name < k.Properties[j].Name
	})

	return nil
}

type Property struct {
	Name string
	Type reflect.Type
}

func (p Property) String() string {
	return fmt.Sprintf("%-15s: %v", p.Name, p.Type)
}

type KeyType int

const (
	StringKeyType KeyType = iota
	IntKeyType
)

func (t KeyType) String() string {
	switch t {
	case StringKeyType:
		return "String"
	case IntKeyType:
		return "Int"
	}
	return "Undefined"
}

/*
type PropertyValue struct {
	Repr []string `datastore:"property_representation"`
}
*/

func GetKinds(ctx context.Context, p *config.Project, names ...string) ([]*Kind, error) {

	setEnv(p)
	//namespace 指定の場合
	cli, err := datastore.NewClient(ctx, p.ProjectID)
	if err != nil {
		return nil, xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	q := datastore.NewQuery("__kind__").KeysOnly()
	keys, err := cli.GetAll(ctx, q, nil)
	if err != nil {
		return nil, xerrors.Errorf("kind GetAll() error: %w", err)
	}

	var kinds []*Kind
	for _, elm := range keys {

		kindName := elm.Name
		var kind Kind
		kind.Name = kindName

		if len(names) != 0 {
			for _, name := range names {
				if kindName == name {
					kinds = append(kinds, &kind)
					break
				}
			}
			if len(kinds) == len(names) {
				break
			}
		} else {
			kinds = append(kinds, &kind)
		}
	}

	if len(names) != 0 {
		if len(kinds) != len(names) {
			return nil, fmt.Errorf("NotFound KindName")
		}
	}

	sort.Slice(kinds, func(i, j int) bool {
		return kinds[i].Name < kinds[j].Name
	})

	for _, kind := range kinds {
		err = kind.setProperties(cli, ctx)
		if err != nil {
			return nil, xerrors.Errorf("Kind[%s] setProperties() error: %w", kind.Name, err)
		}
	}
	return kinds, nil
}

func RemoveKind(ctx context.Context, p *config.Project, name string) error {

	setEnv(p)
	cli, err := datastore.NewClient(ctx, p.ProjectID)
	if err != nil {
		return xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	_, err = cli.RunInTransaction(ctx, func(tx *datastore.Transaction) error {

		fmt.Println("---------------", name)

		q := datastore.NewQuery(name).KeysOnly()
		keys, err := cli.GetAll(ctx, q, nil)
		if err != nil {
			return xerrors.Errorf("GetAll()[%s]: %w", name, err)
		}

		err = tx.DeleteMulti(keys)
		if err != nil {
			return xerrors.Errorf("delete multi[%s]: %w", name, err)
		}

		return nil
	})

	if err != nil {
		return xerrors.Errorf("remove kind error: %w", err)
	}
	return nil
}
