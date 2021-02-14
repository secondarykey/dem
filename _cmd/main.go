package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/secondarykey/dem"
	"golang.org/x/xerrors"
)

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

	switch args[0] {
	case "remove":
		err = dem.Remove()
	case "view":
		err = dem.View()
	default:
		err = fmt.Errorf("Not Support sub command[%s:", args[0])
	}

	if err != nil {
		return xerrors.Errorf("dem error: %w", err)
	}

	return nil
}
