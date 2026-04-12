package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func GenerateInsertQuery(table string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := strings.Split(modelType.Field(i).Tag.Get("db"), ",")[0]
		fmt.Println(dbTag)
		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += " ,"
				placeholders += " ,"
			}
			columns += dbTag
			placeholders += "?"
		}
	}
	return fmt.Sprintf("Insert into %s (%s) values (%s)", table, columns, placeholders)
}

func GetStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
	modeltype := modelValue.Type()
	values := []interface{}{}

	for i := 0; i < modeltype.NumField(); i++ {
		dbTag := strings.Split(modeltype.Field(i).Tag.Get("db"), ",")[0]
		if dbTag != "" && dbTag != "id" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	return values

}

func isValidSortOrder(order string) bool {
	return (order == "asc" || order == "desc")
}

func isValidFiled(filed string, params map[string]string) bool {
	_, ok := params[filed]
	return ok
}

func AddFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"subject":    "subject",
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
	}

	for param, dbFiled := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " and " + dbFiled + "=?"
			args = append(args, value)
		}
	}
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {

		var sortParts []string
		for _, sortParam := range sortParams {
			parts := strings.Split(sortParam, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortOrder(order) || !isValidFiled(field, params) {
				continue
			}
			sortParts = append(sortParts, field+" "+order)

		}
		if len(sortParts) > 0 {
			query += " order by " + strings.Join(sortParts, ", ")
		}
	}

	fmt.Println(query)
	return query, args
}
