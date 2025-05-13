package main

import (
	"fmt"
	"reflect"

	tidbcfg "github.com/pingcap/tidb/pkg/config"
	//lightningcfg "github.com/pingcap/tidb/pkg/lightning/config"
)

func main() {
	tidb := tidbcfg.Config{}
	fmt.Println("TiDB config:")
	for _, param := range List(tidb) {
		fmt.Printf("Type: %s\tKey: %s\n", param.Type, param.Key)
	}
}
func List(input tidbcfg.Config) []Param {
	var r func(reflect.Value, string) []Param
	r = func(reflectV reflect.Value, path string) (result []Param) {
		if reflectV.Type().Name() == "AtomicBool" || reflectV.Type().Name() == "nullableBool" {
			return []Param{{path, "bool"}}
		}
		switch reflectV.Kind() {
		case reflect.Struct:
			for i := 0; i < reflectV.NumField(); i++ {
				//childpath := strings.ToLower(reflectV.Type().Field(i).Name)
				fullKey := path
				paramName := reflectV.Type().Field(i).Tag.Get("toml")
				switch paramName {
				case "":
					// check if the field is anonymous
					if reflectV.Type().Field(i).Anonymous {
						// do nothing
					} else {
						panic(nil)
					}
				case "-":
					continue
				default:
					if path != "" {
						fullKey += "."
					}
					fullKey += paramName
				}

				result = append(result, r(reflectV.Field(i), fullKey)...)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			result = []Param{{path, "number"}}
		case reflect.String:
			result = []Param{{path, "string"}}
		case reflect.Bool:
			result = []Param{{path, "bool"}}
		case reflect.Slice, reflect.Map:
			result = []Param{{path, "json"}}
		default:
			fmt.Println("Unknown type:", reflectV.Kind(), path)
		}
		return result
	}
	return r(reflect.ValueOf(&input).Elem(), "")
}

type Param struct {
	Key  string
	Type string
}
