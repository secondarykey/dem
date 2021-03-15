package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
)

type TempEntity struct {
	Key     *datastore.Key `datastore:"__key__"`
	Value   int
	Buffer1 string
	Buffer2 string
	Buffer3 string
	Buffer4 string
	Buffer5 string
	Buffer6 string
	Buffer7 string
	Buffer8 string
	Buffer9 string
}

func main() {

	err := run()
	if err != nil {
		log.Printf("run() error: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success")
}

func run() error {
	id := "test-endpoint"

	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")
	os.Setenv("DATASTORE_PROJECT_ID", id)
	os.Setenv("DATASTORE_DATASET", id)

	ctx := context.Background()

	cli, err := datastore.NewClient(ctx, id)
	if err != nil {
		return xerrors.Errorf("datastore.NewClient() error: %w", err)
	}

	var keys []*datastore.Key
	var src []*TempEntity

	for i := 0; i < 10; i++ {
		entity := TempEntity{}
		entity.Value = i
		entity.Buffer1 = strings.Repeat("a", 100)
		entity.Buffer2 = strings.Repeat("b", 100)
		entity.Buffer3 = strings.Repeat("x", 100)
		entity.Buffer4 = strings.Repeat("d", 100)
		entity.Buffer5 = strings.Repeat("e", 100)
		entity.Buffer6 = strings.Repeat("t", 100)
		entity.Buffer7 = strings.Repeat("6", 100)
		entity.Buffer8 = strings.Repeat("0", 100)
		entity.Buffer9 = strings.Repeat("l", 100)

		key := datastore.IncompleteKey("TestKind", nil)
		key.Namespace = "test"
		keys = append(keys, key)
		src = append(src, &entity)
	}

	mp, err := cli.PutMulti(ctx, keys, src)
	if err != nil {
		return xerrors.Errorf("PutMulti() error: %w", err)
	}

	for _, elm := range mp {
		fmt.Println("put key", elm.Namespace)
	}

	q := datastore.NewQuery("TestKind").KeysOnly()

	keys, err = cli.GetAll(ctx, q, nil)
	fmt.Println("TestKind len:", len(keys))
	if err != nil {
		return xerrors.Errorf("non namespace GetAll() error: %w", err)
	}

	q = q.Namespace("test")
	keys, err = cli.GetAll(ctx, q, nil)
	if err != nil {
		return xerrors.Errorf("namespace GetAll() error: %w", err)
	}

	fmt.Println("UseNamespace TestKind len:", len(keys))

	return nil
}
