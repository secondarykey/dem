package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/secondarykey/dem/config"
	"golang.org/x/xerrors"
)

func viewMain(w http.ResponseWriter, dto interface{}) error {
	return view(w, dto, "main.tmpl")
}

func view(w http.ResponseWriter, dto interface{}, files ...string) error {

	funcMap := template.FuncMap{
		//"add": func(a, b int) int { return a + b },
	}

	viewFiles := make([]string, 0, len(files)+1)
	viewFiles = append(viewFiles, "layout.tmpl")
	viewFiles = append(viewFiles, files...)

	root := template.New("layout.tmpl").Funcs(funcMap)

	tmpl, err := root.ParseFS(templates, viewFiles...)
	if err != nil {
		return xerrors.Errorf("template.ParseFS() error: %w", err)
	}

	err = tmpl.Execute(w, dto)
	if err != nil {
		return xerrors.Errorf("template.Execute() error: %w", err)
	}
	return nil
}

func viewError(w http.ResponseWriter, msg string, no int, err error) {

	detail := fmt.Sprintf("%+v", err)

	log.Println(err)

	dto := struct {
		Status  int
		Message string
		Detail  string
	}{no, msg, detail}

	w.WriteHeader(no)
	ve := view(w, dto, "error.tmpl")
	if ve != nil {
		log.Printf("viewError() view error: %+v\n", ve)
	}
}

func viewJSON(w http.ResponseWriter, dto interface{}) error {

	res, err := json.Marshal(dto)
	if err != nil {
		return xerrors.Errorf("json.Marshal() error: %w", err)
	}

	w.Header().Set("content-type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		return xerrors.Errorf("writer.Write() error: %w", err)
	}
	return nil
}

func errorJSON(w http.ResponseWriter, msg string, no int, err error) {

	detail := fmt.Sprintf("%+v", err)

	log.Println(detail)

	dto := struct {
		Success bool
		Message string
		Detail  string
	}{false, msg, detail}

	res, err := json.Marshal(dto)
	if err != nil {
		log.Printf("errorJSON() json.Marshal() error:%+v\n", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(no)

	_, err = w.Write(res)
	if err != nil {
		log.Printf("errorJSON() writer.Write() error:%+v\n", err)
		return
	}
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
