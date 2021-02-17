package handler

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed _assets
var embAssets embed.FS
var assets fs.FS

func init() {
	var err error
	assets, err = fs.Sub(embAssets, "_assets")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func registerStatic() error {
	http.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.FS(assets))))
	return nil
}
