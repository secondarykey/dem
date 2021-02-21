package handler

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"

	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

//go:embed _templates
var embTemplates embed.FS
var templates fs.FS

type IndexDto struct {
	Projects   []*config.Project
	Kinds      []*datastore.Kind
	Namespaces []string
	Title      string
	DarkMode   bool
	Current    config.Embed
}

func init() {
	var err error
	templates, err = fs.Sub(embTemplates, "_templates")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func Register() error {

	err := registerStatic()
	if err != nil {
		return xerrors.Errorf("registerStatic() error: %w", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/view/dark/{val}", changeDarkModeHandler)

	r.HandleFunc("/project/remove/{id}", deleteProjectHandler)
	r.HandleFunc("/project/add.json", registerProjectHandler)

	r.HandleFunc("/namespace/change", changeCurrentHandler)
	r.HandleFunc("/entity/limit/change", changeCurrentHandler)

	r.HandleFunc("/entity/view", viewKindHandler)
	r.HandleFunc("/entity/remove", removeEntityHandler)

	r.HandleFunc("/{id}/", viewProjectHandler)

	r.HandleFunc("/", indexHandler)

	http.Handle("/", r)
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	current := config.Embed{}
	current.ID = ""
	current.Kind = ""
	current.Limit = config.DefaultLimit
	current.Namespace = config.DefaultNamespace
	config.SetCurrentEmbed(&current)

	projects, err := config.GetProjects()
	if err != nil {
		viewError(w, "Failed to get project", 500, err)
		return
	}

	dto := IndexDto{}
	dto.Projects = projects
	dto.Kinds = nil
	nss := []string{config.DefaultNamespace}
	dto.Namespaces = nss
	dto.Title = "Select Project"
	dto.DarkMode = config.GetDarkMode()

	dto.Current = current

	err = viewMain(w, dto)
	if err != nil {
		viewError(w, "Failed to write html", 500, err)
	}
}

func changeDarkModeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	darkMode := vars["val"]

	v, err := strconv.ParseBool(darkMode)
	if err != nil {
		viewError(w, "Failed to parse dark mode["+darkMode+"]", 500, err)
		return
	}

	err = config.SetDarkMode(v)
	if err != nil {
		viewError(w, "Failed to set dark mode", 500, err)
		return
	}

	http.Redirect(w, r, "/", 302)
}

func viewProjectHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	current := config.Embed{}
	current.ID = id
	current.Limit = config.DefaultLimit
	current.Namespace = config.DefaultNamespace
	current.Kind = ""
	config.SetCurrentEmbed(&current)

	p := config.SwitchProject(id)
	if p == nil {
		viewError(w, "Failed to switch project", 500, fmt.Errorf("NotFound Project["+id+"]"))
		return
	}

	projects, err := config.GetProjects()
	if err != nil {
		viewError(w, "Failed to get projects", 500, err)
		return
	}

	ctx := r.Context()
	nss, err := datastore.GetNamespaces(ctx)
	if err != nil {
		viewError(w, "Failed to get namespaces", 500, err)
		return
	}

	kinds, err := datastore.GetKinds(ctx)
	if err != nil {
		viewError(w, "Failed to get kinds", 500, err)
		return
	}

	dto := IndexDto{}
	dto.Projects = projects
	dto.Kinds = kinds
	dto.Namespaces = nss
	dto.Title = fmt.Sprintf("%s[%s]", p.Endpoint, p.ProjectID)
	dto.DarkMode = config.GetDarkMode()

	dto.Current = current

	//現在の設定でKindを取得
	err = viewMain(w, dto)
	if err != nil {
		viewError(w, "Failed to write html", 500, err)
	}
}
