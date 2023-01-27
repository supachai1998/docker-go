package helpers

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func SortStruct(slice interface{}, fieldName string, orderTypes ...string) (interface{}, error) {

	orderType := ""
	if orderTypes != nil && len(orderTypes) > 0 {
		orderType = strings.ToLower(orderTypes[0])
	}
	// convert fieldName snek_case to camelCase
	fieldName = strings.ReplaceAll(fieldName, "_", "")

	// slice must be not empty
	if slice == nil || reflect.ValueOf(slice).Len() == 0 {
		return nil, fmt.Errorf("slice must be not empty")
	}

	rv := reflect.ValueOf(slice)
	if rv.Kind().String() != "slice" {
		return nil, fmt.Errorf("slice must be a slice of structs")
	}

	indexOfField := -1
	// nameOfField := ""
	for i := 0; i < rv.Len(); i++ {
		if rv.Index(i).Kind().String() != "struct" {
			return nil, fmt.Errorf("slice must be a slice of structs")
		}
		// get field name
		for j := 0; j < rv.Index(i).NumField(); j++ {
			if strings.ToLower(rv.Index(i).Type().Field(j).Name) == strings.ToLower(fieldName) {
				indexOfField = j
				// nameOfField = rv.Index(i).Type().Field(j).Name
				// fmt.Println("found field name: ", nameOfField, " at index: ", indexOfField, " of struct: ", rv.Index(i).Type().Name(), " in slice")
				break
			}
		}

	}
	// fmt.Println("found name of field: ", nameOfField, " at index: ", indexOfField, " of struct: ", rv.Index(0).Type().Name(), " in slice")
	if indexOfField == -1 {
		return nil, fmt.Errorf("field name not found")
	}

	// sort slice
	switch rv.Index(0).Field(indexOfField).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sort.Slice(slice, func(i, j int) bool {
			if orderType == "desc" {
				return rv.Index(i).Field(indexOfField).Int() > rv.Index(j).Field(indexOfField).Int()
			}
			return rv.Index(i).Field(indexOfField).Int() < rv.Index(j).Field(indexOfField).Int()
		})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sort.Slice(slice, func(i, j int) bool {
			if orderType == "desc" {
				return rv.Index(i).Field(indexOfField).Uint() > rv.Index(j).Field(indexOfField).Uint()
			}
			return rv.Index(i).Field(indexOfField).Uint() < rv.Index(j).Field(indexOfField).Uint()
		})
	case reflect.Float32, reflect.Float64:
		sort.Slice(slice, func(i, j int) bool {
			if orderType == "desc" {
				return rv.Index(i).Field(indexOfField).Float() > rv.Index(j).Field(indexOfField).Float()
			}
			return rv.Index(i).Field(indexOfField).Float() < rv.Index(j).Field(indexOfField).Float()
		})
	case reflect.String:
		sort.Slice(slice, func(i, j int) bool {
			if orderType == "desc" {
				return rv.Index(i).Field(indexOfField).String() > rv.Index(j).Field(indexOfField).String()
			}
			return rv.Index(i).Field(indexOfField).String() < rv.Index(j).Field(indexOfField).String()
		})
		// case pointer
	case reflect.Ptr:
		switch rv.Index(0).Field(indexOfField).Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sort.Slice(slice, func(i, j int) bool {
				if orderType == "desc" {
					return rv.Index(i).Field(indexOfField).Elem().Int() > rv.Index(j).Field(indexOfField).Elem().Int()
				}
				return rv.Index(i).Field(indexOfField).Elem().Int() < rv.Index(j).Field(indexOfField).Elem().Int()
			})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			sort.Slice(slice, func(i, j int) bool {
				if orderType == "desc" {
					return rv.Index(i).Field(indexOfField).Elem().Uint() > rv.Index(j).Field(indexOfField).Elem().Uint()
				}
				return rv.Index(i).Field(indexOfField).Elem().Uint() < rv.Index(j).Field(indexOfField).Elem().Uint()
			})
		case reflect.Float32, reflect.Float64:
			sort.Slice(slice, func(i, j int) bool {
				if orderType == "desc" {
					return rv.Index(i).Field(indexOfField).Elem().Float() > rv.Index(j).Field(indexOfField).Elem().Float()
				}
				return rv.Index(i).Field(indexOfField).Elem().Float() < rv.Index(j).Field(indexOfField).Elem().Float()
			})
		case reflect.String:
			sort.Slice(slice, func(i, j int) bool {
				if orderType == "desc" {
					return rv.Index(i).Field(indexOfField).Elem().String() > rv.Index(j).Field(indexOfField).Elem().String()
				}
				return rv.Index(i).Field(indexOfField).Elem().String() < rv.Index(j).Field(indexOfField).Elem().String()
			})
		default:
			return nil, fmt.Errorf("field type not supported")

		}
	default:
		return nil, fmt.Errorf("field type not supported")
	}

	return slice, nil
}
