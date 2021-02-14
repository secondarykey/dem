package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/secondarykey/dem/config"
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
		fmt.Printf("%T", elm.Value)
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

func GetEntities(ctx context.Context, p *config.Project, name string) ([]*Entity, error) {

	return nil, nil
}
