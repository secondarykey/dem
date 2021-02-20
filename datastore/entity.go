package datastore

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
)

type Entity struct {
	Key    *datastore.Key `datastore:"__key__"`
	Values map[string]interface{}
}

func (e *Entity) LoadKey(k *datastore.Key) error {
	e.Key = k
	return nil
}

func (e *Entity) Load(props []datastore.Property) error {
	e.Values = make(map[string]interface{})
	for _, elm := range props {
		e.Values[elm.Name] = elm.Value
	}
	return nil
}

func (e *Entity) Save() ([]datastore.Property, error) {

	props := make([]datastore.Property, 0, len(e.Values))

	for key, elm := range e.Values {
		p := datastore.Property{}
		p.Name = key
		p.Value = elm
		props = append(props, p)
	}

	return props, nil
}

func (e *Entity) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("=== Key(%d)[%s]", e.Key.ID, e.Key.Name))

	for key, elm := range e.Values {
		line := fmt.Sprintf("%v", elm)
		if len(line) > 100 {
			line = line[0:84] + "..."
		}

		b.WriteString(fmt.Sprintf("\n  %-12s:%v", key, line))
	}

	return b.String()
}

func GetEntities(ctx context.Context, name string) ([]*Entity, error) {

	id, err := setEnvironment()
	if err != nil {
		return nil, xerrors.Errorf("setEnvironment() error:%w", err)
	}

	cli, err := datastore.NewClient(ctx, id)
	if err != nil {
		return nil, xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	q := datastore.NewQuery(name)

	var dst []*Entity
	_, err = cli.GetAll(ctx, q, &dst)
	if err != nil {
		return nil, xerrors.Errorf("GetAll() error: %w", err)
	}

	return dst, nil
}

func RemoveEntity(ctx context.Context, name string, ids []string) error {

	id, err := setEnvironment()
	if err != nil {
		return xerrors.Errorf("setEnvironment() error: %w", err)
	}

	cli, err := datastore.NewClient(ctx, id)
	if err != nil {
		return xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	kinds, err := GetKinds(ctx, name)
	if err != nil {
		return xerrors.Errorf("GetKinds() error: %w", err)
	}

	kind := kinds[0]

	_, err = cli.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		keys := make([]*datastore.Key, len(ids))
		for idx, v := range ids {
			var key *datastore.Key
			switch kind.KeyType {
			case IntKeyType:
				i, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return xerrors.Errorf("strconv.ParseInt() error: %w", err)
				}
				key = datastore.IDKey(name, i, nil)
			case StringKeyType:
				key = datastore.NameKey(name, v, nil)
			}
			keys[idx] = key
		}

		err := tx.DeleteMulti(keys)
		if err != nil {
			return xerrors.Errorf("transaction remove all error: %w", err)
		}

		return nil
	})

	if err != nil {
		return xerrors.Errorf("remove transaction error: %w", err)
	}

	return nil
}
