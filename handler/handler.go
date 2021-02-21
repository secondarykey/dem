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
	ID         string
	DarkMode   bool
	Limit      int
	Namespace  string
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

	r.HandleFunc("/project/delete.json", deleteProjectHandler)
	r.HandleFunc("/project/add.json", registerProjectHandler)

	r.HandleFunc("/namespace/{ns}", namespaceHandler)

	r.HandleFunc("/entity/view/{kind}", viewKindHandler)
	r.HandleFunc("/entity/limit/{limit}", changeLimitHandler)
	r.HandleFunc("/entity/remove/{kind}", removeEntityHandler)

	r.HandleFunc("/{id}/", viewProjectHandler)

	r.HandleFunc("/", indexHandler)

	http.Handle("/", r)
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	projects, err := config.GetProjects()
	if err != nil {
		log.Println(err)
		return
	}

	dto := IndexDto{}
	dto.Projects = projects
	dto.Kinds = nil
	dto.Title = "Select Project"
	dto.ID = "empty"
	dto.DarkMode = config.GetDarkMode()
	dto.Limit = config.GetLimit()

	nss := []string{config.DefaultNamespace}
	dto.Namespaces = nss
	dto.Namespace = config.DefaultNamespace

	err = viewMain(w, dto)
	if err != nil {
		log.Println(err)
	}
}

func changeDarkModeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	darkMode := vars["val"]

	v, err := strconv.ParseBool(darkMode)
	if err != nil {
		log.Println(err)
		return
	}

	err = config.SetDarkMode(v)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/", 302)
}

func viewProjectHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	p := config.SwitchProject(id)
	if p == nil {
		log.Println("NotFound")
		return
	}

	projects, err := config.GetProjects()
	if err != nil {
		log.Println(err)
		return
	}

	ctx := r.Context()
	nss, err := datastore.GetNamespaces(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	kinds, err := datastore.GetKinds(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	dto := IndexDto{}
	dto.Projects = projects
	dto.Kinds = kinds
	dto.Title = fmt.Sprintf("%s[%s]", p.Endpoint, p.ProjectID)
	dto.ID = p.ID
	dto.DarkMode = config.GetDarkMode()
	dto.Limit = config.GetLimit()
	dto.Namespaces = nss
	dto.Namespace = config.DefaultNamespace

	config.SetNamespace(dto.Namespace)

	//現在の設定でKindを取得
	err = viewMain(w, dto)
	if err != nil {
		log.Println(err)
	}
}
