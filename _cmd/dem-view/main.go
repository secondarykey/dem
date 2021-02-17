package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/secondarykey/dem"
	"github.com/secondarykey/dem/config"

	"golang.org/x/xerrors"
)

var port int
var path string

func init() {
	flag.IntVar(&port, "p", 8088, "Viewer port number.")
	flag.StringVar(&path, "c", "$HOME/.dem.gob", "Config file path")
}

func main() {

	flag.Parse()

	err := run()
	if err != nil {
		fmt.Printf("dem error:\n%+v\n", err)
		os.Exit(1)
	}
	fmt.Println("Bye!")
}

func run() error {
	err := dem.Listen(
		config.SetMyPort(port),
		config.SetConfigFile(path))
	if err != nil {
		return xerrors.Errorf("dem.Listen() error: %w", err)
	}
	return nil
}
