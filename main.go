package main

import (
	"fmt"
	"reflect"

	tidbcfg "github.com/pingcap/tidb/pkg/config"
	lightningcfg "github.com/pingcap/tidb/pkg/lightning/config"
	pdcfg "github.com/tikv/pd/server/config"
)

func main() {
	test := TestConfig{}
	for _, param := range List(test) {
		fmt.Printf("Type: %s\tKey: %s\n", param.Type, param.Key)
	}
	tidb := tidbcfg.Config{}
	fmt.Println("TiDB config:")
	for _, param := range List(tidb) {
		fmt.Printf("Type: %s\tKey: %s\n", param.Type, param.Key)
	}

	lightning := lightningcfg.Config{}
	fmt.Println("\nLightning config:")
	for _, param := range List(lightning) {
		fmt.Printf("Type: %s\tKey: %s\n", param.Type, param.Key)
	}
	pd := pdcfg.NewConfig()
	fmt.Println("\nPD config:")
	for _, param := range List(pd) {
		fmt.Printf("Type: %s\tKey: %s\n", param.Type, param.Key)
	}
}
func List(input interface{}) []Param {
	var r func(reflect.Value, string) []Param
	r = func(reflectV reflect.Value, path string) (result []Param) {
		switch reflectV.Type().Name() {
		case "AtomicBool", "nullableBool":
			return []Param{{path, "bool"}}
		case "Int64":
			return []Param{{path, "number"}}
		}
		switch reflectV.Kind() {
		case reflect.Struct:
			for i := 0; i < reflectV.NumField(); i++ {
				fullKey := path
				paramName := reflectV.Type().Field(i).Tag.Get("toml")
				switch paramName {
				case "":
					if reflectV.Type().Field(i).Anonymous {
						// do nothing to passthrough to the child field
					} else {
						fmt.Println("ignore no tag field: "+path+" "+reflectV.Type().Field(i).Name)
						continue
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
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = []Param{{path, "int"}}
		case reflect.Float32, reflect.Float64:
			result = []Param{{path, "float"}}
		case reflect.String:
			result = []Param{{path, "string"}}
		case reflect.Bool:
			result = []Param{{path, "bool"}}
		case reflect.Slice, reflect.Map:
			result = []Param{{path, "json"}}
		case reflect.Interface:
			result = append(result, r(reflectV.Elem(), path)...)
		case reflect.Ptr:
			result = append(result, r(reflect.Zero(reflectV.Type().Elem()), path)...)

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

type TestConfig struct {
	Int int     `toml:"int"`
	Str string  `toml:"str"`
	Ptr *string `toml:"ptr"`
}
