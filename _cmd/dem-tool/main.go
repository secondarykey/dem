package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/secondarykey/dem"
	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

var (
	port      int
	host      string
	projectID string
	ns        string
)

func init() {
	flag.IntVar(&port, "p", 8081, "Datastore Emulator port number.")
	flag.StringVar(&host, "h", "localhost", "Datastore Emulator Host.")
	flag.StringVar(&projectID, "project", "project", "Datastore Emulator ProjectID.")
	flag.StringVar(&ns, "n", "", "Namespace")
}

func main() {

	if err := run(); err != nil {
		fmt.Printf("dem run error:%+v\n", err)
		os.Exit(1)
	}
	fmt.Println("Success")
}

func run() error {

	flag.Parse()
	args := flag.Args()

	if len(args) <= 0 {
		return xerrors.Errorf("run() arguments error")
	}
	var err error

	kinds := make([]string, 0)
	if len(args) >= 1 {
		kinds = args[1:]
	}

	err = dem.SetConsoleOptions(
		config.SetEmulatorPort(port),
		config.SetHost(host),
		config.SetProjectID(projectID),
		config.SetNamespace(ns))
	if err != nil {
		return xerrors.Errorf("dem.SetOptions error: %w", err)
	}

	switch args[0] {
	case "remove":
		err = dem.RemoveEntity(kinds...)
	case "kind":
		err = dem.ViewKind(kinds...)
	case "entity":
		err = dem.ViewEntity(kinds...)
	default:
		err = fmt.Errorf("Not Support sub command[%s:", args[0])
	}

	if err != nil {
		return xerrors.Errorf("dem error: %w", err)
	}

	return nil
}
