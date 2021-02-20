package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/secondarykey/dem/datastore"
)

func removeEntityHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	buf := r.FormValue("ids")
	var ids []string

	err := json.Unmarshal([]byte(buf), &ids)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	kind := vars["kind"]

	err = datastore.RemoveEntity(r.Context(), kind, ids)
	if err != nil {
		log.Println(err)
		return
	}

	dto := struct {
		Success bool
	}{true}

	err = viewJSON(w, dto)
	if err != nil {
		log.Println(err)
	}

}
