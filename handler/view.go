package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

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

func createCurrent(r *http.Request) *config.Embed {

	current := config.Embed{}

	r.ParseForm()
	current.ID = r.FormValue("ID")
	current.Kind = r.FormValue("kind")
	current.Namespace = r.FormValue("namespace")

	lmtBuf := r.FormValue("limit")
	v, err := strconv.Atoi(lmtBuf)
	if err != nil {
		log.Println("limit parse error: %+v", v)
		v = config.DefaultLimit
	}

	current.Limit = v
	return &current
}
