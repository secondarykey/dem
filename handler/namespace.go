package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/secondarykey/dem/config"
)

func namespaceHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ns := vars["ns"]

	config.SetNamespace(ns)

	dto := struct {
		Success bool
	}{true}

	err := viewJSON(w, dto)
	if err != nil {
		log.Println(err)
	}
}
