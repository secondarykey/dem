package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/secondarykey/dem/config"
	"github.com/secondarykey/dem/datastore"
)

type Entity struct {
	Key    string
	Types  []ViewType
	Values []Value
}

type Value struct {
	View string
	Real string
}

type ViewType int

const (
	NormalType   ViewType = ViewType(0)
	OmittedType           = ViewType(10)
	ExpandType            = ViewType(20)
	SliceType             = ViewType(30)
	DownloadType          = ViewType(40)
	ErrorType             = ViewType(-1)
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

	props, obj := convertViewEntity(kind, entity)

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

	props, values := convertViewEntities(kinds[0], entities)

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

//
// convertViewEntiry is displayed dialog data
//
func convertViewEntity(kind *datastore.Kind, entity *datastore.Entity) ([]string, *Entity) {

	datum := Entity{}
	key := entity.Key
	if kind.KeyType == datastore.StringKeyType {
		datum.Key = key.Name
	} else {
		datum.Key = fmt.Sprintf("%d", key.ID)
	}

	datum.Values = make([]Value, len(kind.Properties))
	datum.Types = make([]ViewType, len(kind.Properties))
	header := make([]string, len(kind.Properties))
	for idx, prop := range kind.Properties {

		header[idx] = prop.Name
		val, ok := entity.Values[prop.Name]
		v := ""
		r := ""
		t := NormalType
		if !ok {
			v = "Mismatch " + prop.Name
			t = ErrorType
		} else {
			switch nv := val.(type) {
			case []uint8:
				v = fmt.Sprintf("byte length %d", len(nv))
				t = DownloadType
			case time.Time:
				v = fmt.Sprintf("%v", nv)
			case string:
				cut := cutData(nv)
				if cut != nv {
					t = OmittedType
					r = nv
					v = cut
				} else {
					t = NormalType
					r = nv
					v = nv
				}
			case datastore.Entity:

				v = fmt.Sprintf("%v", nv)

			default:
				v = fmt.Sprintf("%v", nv)
			}
		}

		datum.Types[idx] = t

		setV := Value{}
		setV.View = v
		setV.Real = r

		datum.Values[idx] = setV
	}

	return header, &datum
}

//
// convertViewEntitie() is displayed table data
//
func convertViewEntities(kind *datastore.Kind, entities []*datastore.Entity) ([]string, []*Entity) {

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

		vals := make([]Value, len(kind.Properties)+1)
		keyVal := Value{}
		keyVal.View = kv
		vals[0] = keyVal

		for jdx, prop := range kind.Properties {
			v, ok := entity.Values[prop.Name]
			if !ok {
				v = "Mismatch," + prop.Name
			}
			viewVal := Value{}
			viewVal.View = cutData(v)
			vals[jdx+1] = viewVal
		}

		datum.Values = vals
		data[idx] = &datum
	}

	return props, data
}

func cutData(v interface{}) string {

	switch val := v.(type) {
	case string:

		n := false
		if idx := strings.Index(val, "\n"); idx != -1 {
			val = val[:idx-1]
			n = true
		}

		r := []rune(val)
		if len(r) > 15 {
			return string(r[0:13]) + "..."
		}

		if n {
			return val + "..."
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
