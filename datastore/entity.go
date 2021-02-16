package datastore

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/secondarykey/dem/config"
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
		//TODO １度は型を作成する
		//reflect.Typeof()
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
		b.WriteString(fmt.Sprintf("\n  %-12s:%v", key, elm))
	}

	return b.String()
}

func GetEntities(ctx context.Context, p *config.Project, name string) ([]*Entity, error) {

	cli, err := datastore.NewClient(ctx, p.ProjectID)
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
