package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/secondarykey/dem/config"
)

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	err := config.DeleteProject(id)
	if err != nil {
		viewError(w, "Faild to delete project.", 500, err)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func registerProjectHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	endpoint := r.FormValue("endpoint")
	projectID := r.FormValue("projectid")

	pro := config.NewProject(endpoint, projectID)

	err := config.AddProject(pro)
	if err != nil {
		errorJSON(w, "Failed to add project", 500, err)
		return
	}

	dto := struct {
		Success  bool
		Redirect string
	}{true, fmt.Sprintf("/%s/", pro.ID)}

	err = viewJSON(w, dto)
	if err != nil {
		errorJSON(w, "Failed to view json", 500, err)
	}
}
