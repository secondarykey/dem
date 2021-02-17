package handler

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"

	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

//go:embed _templates
var embTemplates embed.FS
var templates fs.FS

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

	r.HandleFunc("/project/delete.json", deleteProjectHandler)
	r.HandleFunc("/project/add.json", registerProjectHandler)

	r.HandleFunc("/namespace/change", namespaceHandler)

	r.HandleFunc("/{id}/", viewProjectHandler)
	r.HandleFunc("/{id}/{kind}/", viewKindHandler)

	r.HandleFunc("/", indexHandler)

	http.Handle("/", r)
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	projects, err := config.GetProjects()
	if err != nil {
		log.Println(err)
	}

	dto := struct {
		Projects []*config.Project
		Kinds    []*datastore.Kind
		Title    string
		ID       string
	}{projects, nil, "Select Project", "empty"}

	err = viewMain(w, dto)
	if err != nil {
		log.Println(err)
	}
}

func viewProjectHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	p := config.GetProject(id)
	if p == nil {
		log.Println("NotFound")
		return
	}

	projects, err := config.GetProjects()
	if err != nil {
		log.Println(err)
		return
	}

	kinds, err := datastore.GetKinds(context.Background(), p)
	if err != nil {
		log.Println(err)
		return
	}

	dto := struct {
		Projects []*config.Project
		Kinds    []*datastore.Kind
		Title    string
		ID       string
	}{projects, kinds, fmt.Sprintf("%s[%s]", p.Endpoint, p.ProjectID), p.ID}

	//現在の設定でKindを取得
	err = viewMain(w, dto)
	if err != nil {
		log.Println(err)
	}
}

func viewMain(w http.ResponseWriter, dto interface{}) error {
	return view(w, dto, "layout.tmpl")
}

func view(w http.ResponseWriter, dto interface{}, files ...string) error {

	funcMap := template.FuncMap{
		//"add": func(a, b int) int { return a + b },
	}

	root := template.New("layout.tmpl").Funcs(funcMap)

	tmpl, err := root.ParseFS(templates, files...)
	if err != nil {
		return xerrors.Errorf("template.ParseFS() error: %w", err)
	}

	err = tmpl.Execute(w, dto)
	if err != nil {
		return xerrors.Errorf("template.Execute() error: %w", err)
	}
	return nil
}

func viewJSON(w http.ResponseWriter, dto interface{}) error {
	w.Header().Set("content-type", "application/json")
	res, err := json.Marshal(dto)
	if err != nil {
		return xerrors.Errorf("json.Marshal() error: %w", err)
	}

	_, err = w.Write(res)
	if err != nil {
		return xerrors.Errorf("writer.Write() error: %w", err)
	}
	return nil
}
