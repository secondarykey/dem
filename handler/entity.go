package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/secondarykey/dem/config"
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

func changeLimitHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	limit := vars["limit"]
	v, err := strconv.Atoi(limit)
	if err != nil {
		log.Println(err)
		return
	}

	config.SetLimit(v)
	dto := struct {
		Success bool
	}{true}

	err = viewJSON(w, dto)
	if err != nil {
		log.Println(err)
	}
}
