package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	http.HandleFunc("/project/view.json", viewProjectHandler)
	http.HandleFunc("/project/delete.json", deleteProjectHandler)
	http.HandleFunc("/project/add.json", registerProjectHandler)
	http.HandleFunc("/{id}/", viewKindHandler)
	http.HandleFunc("/{id}/{kind}/", viewEntityHandler)
	http.Handle("/", r)
}
