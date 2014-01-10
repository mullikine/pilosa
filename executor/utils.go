package executor

import (
	"io/ioutil"
	"pilosa/query"

	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto"
)

func GetMacro(file_name string, filter string) interface{} {

	file_data, err := ioutil.ReadFile(file_name)
	if err != nil {
		spew.Dump(err)
	}
	s := string(file_data[:])

	js := "query_list = (function (filter){" + s + "})('" + filter + "');"

	Otto := otto.New()
	Otto.Run(js)
	query_objects, err := Otto.Get("query_list")

	query_list_interface, err := query_objects.Export()
	if err != nil {
		spew.Dump(err)
	}

	var query_list query.QueryList

	// ql is []interface{}
	switch ql := query_list_interface.(type) {
	case []interface{}:
		spew.Dump("INTERFACE", ql)
		// q is map[string]interface{}
		for i, _ := range ql {
			q := ql[i].(map[string]interface{})
			spew.Dump(ql[i])
			spew.Dump(q)
			spew.Dump(q["label"].(string))
			spew.Dump(q["pql"].(string))
			query_list = append(query_list, query.QueryListItem{Label: q["label"].(string), PQL: q["pql"].(string)})
		}
	default:
		spew.Dump("DEFAULT")
	}

	return query_list
}
