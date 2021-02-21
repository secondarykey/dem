package handler

import (
	"encoding/json"
	"log"
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
		log.Println(err)
		return
	}

	err = datastore.RemoveEntity(r.Context(), c.Kind, ids)
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

func changeCurrentHandler(w http.ResponseWriter, r *http.Request) {

	e := createCurrent(r)
	config.SetCurrent(e)

	dto := struct {
		Success bool
	}{true}

	err := viewJSON(w, dto)
	if err != nil {
		log.Println(err)
	}
}
