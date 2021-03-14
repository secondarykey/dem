package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"
)

func getEntityHandler(w http.ResponseWriter, r *http.Request) {

	e := createCurrent(r)
	config.SetCurrentEmbed(e)
	kinds, err := datastore.GetKinds(r.Context(), e.Kind)
	if err != nil {
		errorJSON(w, "Failed to get Kind["+e.Kind+"]", 500, err)
		return
	}

	kind := kinds[0]

	key := r.FormValue("key")
	strKey := ""
	intKey := 0

	if kind.KeyType == datastore.StringKeyType {
		strKey = key
	} else {
		intKey, err = strconv.Atoi(key)
		if err != nil {
			errorJSON(w, "Failed to parse key["+key+"]", 500, err)
			return
		}
	}

	entity, err := datastore.GetEntity(r.Context(), e.Kind, strKey, int64(intKey))
	if err != nil {
		errorJSON(w, "Failed to get entity["+key+"]", 500, err)
		return
	}

	props, obj := convertDisplayData(kind, entity)

	dto := struct {
		Success bool
		Header  []string
		Entity  *Entity
	}{true, props, obj}

	err = viewJSON(w, dto)
	if err != nil {
		errorJSON(w, "Failed to view json", 500, err)
	}
}

func viewEntitiesHandler(w http.ResponseWriter, r *http.Request) {

	e := createCurrent(r)
	config.SetCurrentEmbed(e)

	cursor := r.FormValue("cursor")
	v := r.FormValue("first")
	first, err := strconv.ParseBool(v)
	if err != nil || first {
		cursor = ""
	}

	kinds, err := datastore.GetKinds(r.Context(), e.Kind)
	if err != nil {
		errorJSON(w, "Failed to get Kind["+e.Kind+"]", 500, err)
		return
	}

	entities, cur, err := datastore.GetEntities(r.Context(), e.Kind, e.Limit, cursor, e.Namespace)
	if err != nil {
		errorJSON(w, "Failed to select entity", 500, err)
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
		errorJSON(w, "Failed to view json", 500, err)
	}
}

func convertDisplayData(kind *datastore.Kind, entity *datastore.Entity) ([]string, *Entity) {

	datum := Entity{}
	key := entity.Key
	if kind.KeyType == datastore.StringKeyType {
		datum.Key = key.Name
	} else {
		datum.Key = fmt.Sprintf("%d", key.ID)
	}

	datum.Values = make([]string, len(kind.Properties))
	header := make([]string, len(kind.Properties))
	for idx, prop := range kind.Properties {

		header[idx] = prop.Name
		val, ok := entity.Values[prop.Name]
		v := ""
		if !ok {
			v = "Mismatch " + prop.Name
		} else {
			switch nv := val.(type) {
			case []uint8:
				v = fmt.Sprintf("byte length %d", len(nv))
			default:
				v = fmt.Sprintf("%v", nv)
			}
			fmt.Printf("%s = [%T]\n", prop.Name, val)
		}
		datum.Values[idx] = v
	}

	return header, &datum
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
