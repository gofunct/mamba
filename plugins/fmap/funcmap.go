package fmap

import (
	"encoding/json"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/mamba/plugins/fmap/protofunc"
	"github.com/huandu/xstrings"
	"github.com/spf13/viper"
	"reflect"
	"strings"
	"text/template"
)

func init() {
	for k, v := range sprig.GenericFuncMap() {
		DefaultFMap[k] = v
	}
	for k, v := range protofunc.DefaultProtoFmap {
		DefaultFMap[k] = v
	}
}

var (
	DefaultFMap = template.FuncMap{
		"string": func(i interface {
			String() string
		}) string {
			return i.String()
		},
		"json": func(v interface{}) string {
			a, err := json.Marshal(v)
			if err != nil {
				return err.Error()
			}
			return string(a)
		},
		"prettyjson": func(v interface{}) string {
			a, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return err.Error()
			}
			return string(a)
		},
		"splitArray": func(sep string, s string) []interface{} {
			var r []interface{}
			t := strings.Split(s, sep)
			for i := range t {
				if t[i] != "" {
					r = append(r, t[i])
				}
			}
			return r
		},
		"first": func(a []string) string {
			return a[0]
		},
		"last": func(a []string) string {
			return a[len(a)-1]
		},
		"concat": func(a string, b ...string) string {
			return strings.Join(append([]string{a}, b...), "")
		},
		"join": func(sep string, a ...string) string {
			return strings.Join(a, sep)
		},
		"upperFirst": func(s string) string {
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"lowerFirst": func(s string) string {
			return strings.ToLower(s[:1]) + s[1:]
		},
		"camelCase": func(s string) string {
			if len(s) > 1 {
				return xstrings.ToCamelCase(s)
			}

			return strings.ToUpper(s[:1])
		},
		"lowerCamelCase": func(s string) string {
			if len(s) > 1 {
				s = xstrings.ToCamelCase(s)
			}

			return strings.ToLower(s[:1]) + s[1:]
		},
		"kebabCase": func(s string) string {
			return strings.Replace(xstrings.ToSnakeCase(s), "_", "-", -1)
		},
		"contains": func(sub, s string) bool {
			return strings.Contains(s, sub)
		},
		"trimstr": func(cutset, s string) string {
			return strings.Trim(s, cutset)
		},
		"index": func(array interface{}, i int32) interface{} {
			slice := reflect.ValueOf(array)
			if slice.Kind() != reflect.Slice {
				panic("Error in index(): given a non-slice type")
			}
			if i < 0 || int(i) >= slice.Len() {
				panic("Error in index(): index out of bounds")
			}
			return slice.Index(int(i)).Interface()
		},
		"add": func(a int, b int) int {
			return a + b
		},
		"subtract": func(a int, b int) int {
			return a - b
		},
		"multiply": func(a int, b int) int {
			return a * b
		},
		"divide": func(a int, b int) int {
			if b == 0 {
				panic("psssst ... little help here ... you cannot divide by 0")
			}
			return a / b
		},
		/////////////////////////////////////////

		"snakeCase": xstrings.ToSnakeCase,

		"set":                      Set,
		"get":                      Get,
		"allSettings":              AllSettings,
		"allKeys":                  AllKeys,
		"configFileUsed":           ConfigFileUsed,
		"supportedExts":            SupportedExts,
		"isSet":                    IsSet,
		"supportedRemoteProviders": SupportedRemoteProviders,
	}
)

func Set(key string, o interface{}) string {
	viper.Set(key, o)
	return ""
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

func AllKeys() []string {
	return viper.AllKeys()
}

func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}
func SupportedExts() []string {
	return viper.SupportedExts
}
func IsSet(key string) bool {
	return viper.IsSet(key)
}

func SupportedRemoteProviders() []string {
	return viper.SupportedRemoteProviders
}

func ReplaceDict(src string, dict map[string]interface{}) string {
	for old, v := range dict {
		new, ok := v.(string)
		if !ok {
			continue
		}
		src = strings.Replace(src, old, new, -1)
	}
	return src
}

func LowerGoNormalize(s string) string {
	fmtd := xstrings.ToCamelCase(s)
	fmtd = xstrings.FirstRuneToLower(fmtd)
	return FormatID(s, fmtd)
}

func GoNormalize(s string) string {
	fmtd := xstrings.ToCamelCase(s)
	return FormatID(s, fmtd)
}

func FormatID(base string, formatted string) string {
	if formatted == "" {
		return formatted
	}
	switch {
	case base == "id":
		// id -> ID
		return "ID"
	case strings.HasPrefix(base, "id_"):
		// id_some -> IDSome
		return "ID" + formatted[2:]
	case strings.HasSuffix(base, "_id"):
		// some_id -> SomeID
		return formatted[:len(formatted)-2] + "ID"
	case strings.HasSuffix(base, "_ids"):
		// some_ids -> SomeIDs
		return formatted[:len(formatted)-3] + "IDs"
	}
	return formatted
}

func GetPackageTypeName(s string) string {
	if strings.Contains(s, ".") {
		return strings.Split(s, ".")[1]
	}
	return ""
}
