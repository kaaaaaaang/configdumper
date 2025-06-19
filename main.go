package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type UnmarshalText interface {
	UnmarshalText([]byte) error
}

type UnmarshalTOML interface {
	UnmarshalTOML(any) error
}

func main() {
	// get tikv config info from stdin
	var tikvConfigInfo []byte
	if input, err := os.ReadFile("/dev/stdin"); err != nil {
		panic(err)
	} else {
		tikvConfigInfo = input
	}
	tikvJSON, err := json.MarshalIndent(ParseTikvConfigInfo(tikvConfigInfo), "", "    ")
	if err != nil {
		panic(err)
	}
	os.WriteFile("/dev/stdout", tikvJSON, 0644)
}

func ParseTomlConfig(input interface{}) []Param {
	var r func(reflect.Value, string) []Param
	r = func(reflectV reflect.Value, path string) (result []Param) {
		switch reflectV.Type().Name() {
		// common
		case "AtomicBool", "nullableBool":
			return []Param{{Key: path, Type: "bool"}}
		case "Int64":
			return []Param{{Key: path, Type: "number"}}
		case "Duration":
			return []Param{{Key: path, Type: "string"}}
		case "ByteSize":
			return []Param{{Key: path, Type: "string"}}
		// PD
		case "RedactInfoLogType":
			return []Param{{Key: path, Type: "unknown"}}
		// lightning
		case "MaxError":
			// it is a struct, do nothing
		case "CheckpointKeepStrategy":
			return []Param{{Key: path, Type: "unknown"}}
		case "StringOrStringSlice":
			return []Param{{Key: path, Type: "json"}}
		case "DuplicateResolutionAlgorithm":
			return []Param{{Key: path, Type: "string"}}
		case "CompressionType":
			return []Param{{Key: path, Type: "string"}}
		case "PostOpLevel":
			// bool or string
			return []Param{{Key: path, Type: "unknown"}}
		default:
			var unmarshalTextType = reflect.TypeOf((*UnmarshalText)(nil)).Elem()
			var unmarshalTOMLType = reflect.TypeOf((*UnmarshalTOML)(nil)).Elem()
			if reflectV.Type().Implements(unmarshalTextType) ||
				(reflect.PointerTo(reflectV.Type()).Implements(unmarshalTextType)) {
				panic(fmt.Sprintf("%s  <-- Implements UnmarshalText, please udate code to handler it", reflectV.Type().Name()))
			} else if reflectV.Type().Implements(unmarshalTOMLType) ||
				(reflect.PointerTo(reflectV.Type()).Implements(unmarshalTOMLType)) {
				panic(fmt.Sprintf("%s  <-- Implements unmarshalTOMLType, please udate code to handler it", reflectV.Type().Name()))
			}
		}
		switch reflectV.Kind() {
		case reflect.Struct:
			for i := 0; i < reflectV.NumField(); i++ {
				fullKey := path
				tag := reflectV.Type().Field(i).Tag.Get("toml")
				paramName := strings.Split(tag, ",")[0]
				switch paramName {
				case "":
					if reflectV.Type().Field(i).Anonymous {
						// do nothing to passthrough to the child field
					} else {
						fmt.Println("ignore no tag field: " + path + " " + reflectV.Type().Field(i).Name)
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
			result = []Param{{Key: path, Type: "int"}}
		case reflect.Float32, reflect.Float64:
			result = []Param{{Key: path, Type: "float"}}
		case reflect.String:
			result = []Param{{Key: path, Type: "string"}}
		case reflect.Bool:
			result = []Param{{Key: path, Type: "bool"}}
		case reflect.Slice, reflect.Map:
			result = []Param{{Key: path, Type: "json"}}
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
	Key          string `json:"key"`
	Type         string `json:"type"`
	DefaultValue any    `json:"default_value,omitempty"`
}

type tikvConfigInfo struct {
	Parameters []struct {
		Name         string `json:"Name"`
		DefaultValue any    `json:"DefaultValue"`
	} `json:"Parameters"`
}

func ParseTikvConfigInfo(input []byte) []Param {
	var configInfo tikvConfigInfo
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	if err := decoder.Decode(&configInfo); err != nil {
		panic(fmt.Sprintf("Failed to decode Tikv config info: %v", err))
	}

	var params []Param
	for _, param := range configInfo.Parameters {
		ttype := "unknown" // Default type is unknown, as we don't have type information in the input
		if param.DefaultValue != nil {
			switch param.DefaultValue.(type) {
			case string:
				ttype = "string"
			case json.Number:
				if strings.Contains(param.DefaultValue.(json.Number).String(), ".") {
					ttype = "float"
				} else {
					ttype = "int"
				}
			case bool:
				ttype = "bool"
			case []interface{}, map[string]interface{}:
				ttype = "json"
			default:
				panic(fmt.Sprintf("Unknown type for parameter %s: %T", param.Name, param.DefaultValue))
			}
		}
		params = append(params, Param{
			Key:          param.Name,
			Type:         ttype,
			DefaultValue: param.DefaultValue,
		})
	}
	return params
}
