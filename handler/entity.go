package handler

import (
	"encoding/json"
	"net/http"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"
)

func removeEntityHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	c := createCurrent(r)

	buf := r.FormValue("ids")
	var ids []string

	err := json.Unmarshal([]byte(buf), &ids)
	if err != nil {
		errorJSON(w, "Failed to unmarshal json", 500, err)
		return
	}

	err = datastore.RemoveEntity(r.Context(), c.Kind, ids)
	if err != nil {
		errorJSON(w, "Failed to remove entities", 500, err)
		return
	}

	dto := struct {
		Success bool
	}{true}

	err = viewJSON(w, dto)
	if err != nil {
		errorJSON(w, "Failed to write json", 500, err)
	}
}

func changeCurrentHandler(w http.ResponseWriter, r *http.Request) {

	e := createCurrent(r)
	config.SetCurrentEmbed(e)

	dto := struct {
		Success bool
	}{true}

	err := viewJSON(w, dto)
	if err != nil {
		errorJSON(w, "Failed to write json", 500, err)
	}
}
