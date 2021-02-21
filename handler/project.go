package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"
)

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
}

func registerProjectHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	endpoint := r.FormValue("endpoint")
	projectID := r.FormValue("projectid")

	pro := config.NewProject(endpoint, projectID)

	err := config.AddProject(pro)
	if err != nil {
		log.Println(err)
		return
	}

	dto := struct {
		Success  bool
		Redirect string
	}{true, fmt.Sprintf("/%s/", pro.ID)}

	err = viewJSON(w, dto)
	if err != nil {
		log.Println(err)
	}

}

func viewKindHandler(w http.ResponseWriter, r *http.Request) {

	e := createCurrent(r)
	config.SetCurrent(e)

	cursor := r.FormValue("cursor")
	v := r.FormValue("first")
	first, err := strconv.ParseBool(v)
	if err != nil || first {
		cursor = ""
	}

	kinds, err := datastore.GetKinds(r.Context(), e.Kind)
	if err != nil {
		log.Println(err)
		return
	}

	entities, cur, err := datastore.GetEntities(r.Context(), e.Kind, e.Limit, cursor, e.Namespace)
	if err != nil {
		log.Println(err)
		return
	}

	props, values := transformData(kinds[0], entities)

	dto := struct {
		Success bool
		Header  []string
		Data    []*Entity
		Next    string
	}{true, props, values, cur}

	err = viewJSON(w, dto)
	if err != nil {
		log.Println(err)
	}
}

func transformData(kind *datastore.Kind, entities []*datastore.Entity) ([]string, []*Entity) {

	props := make([]string, len(kind.Properties)+1)
	props[0] = "Key"
	for idx, prop := range kind.Properties {
		props[idx+1] = prop.Name
	}

	data := make([]*Entity, len(entities))
	for idx, entity := range entities {
		datum := Entity{}
		key := entity.Key
		kv := ""
		if kind.KeyType == datastore.StringKeyType {
			datum.Key = key.Name
			kv = cutData(key.Name)
		} else {
			datum.Key = fmt.Sprintf("%d", key.ID)
			kv = datum.Key
		}

		vals := make([]string, len(kind.Properties)+1)
		vals[0] = kv

		for jdx, prop := range kind.Properties {
			v, ok := entity.Values[prop.Name]
			if !ok {
				v = "Mismatch," + prop.Name
			}
			vals[jdx+1] = cutData(v)
		}

		datum.Values = vals
		data[idx] = &datum
	}

	return props, data
}

func cutData(v interface{}) string {

	switch val := v.(type) {
	case string:
		r := []rune(val)
		if len(r) > 15 {
			return string(r[0:13]) + "..."
		}
		return val
	default:
		str := fmt.Sprintf("%v", v)
		if len(str) > 21 {
			str = str[0:21] + "..."
		}
		return str
	}
	return ""
}

type Entity struct {
	Key    string
	Values []string
}
