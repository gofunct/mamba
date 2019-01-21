package fmap

import (
	"encoding/json"
	"github.com/huandu/xstrings"
	"github.com/spf13/viper"
	"reflect"
	"strings"
	"text/template"
	"time"
)

var ProtoHelpersFuncMap = template.FuncMap{
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

	"snakeCase":                    xstrings.ToSnakeCase,
	"getProtoFile":                 getProtoFile,
	"getMessageType":               getMessageType,
	"getEnumValue":                 getEnumValue,
	"isFieldMessage":               isFieldMessage,
	"isFieldMessageTimeStamp":      isFieldMessageTimeStamp,
	"isFieldRepeated":              isFieldRepeated,
	"haskellType":                  haskellType,
	"goType":                       goType,
	"goZeroValue":                  goZeroValue,
	"goTypeWithPackage":            goTypeWithPackage,
	"goTypeWithGoPackage":          goTypeWithGoPackage,
	"jsType":                       jsType,
	"jsSuffixReserved":             jsSuffixReservedKeyword,
	"namespacedFlowType":           namespacedFlowType,
	"httpVerb":                     httpVerb,
	"httpPath":                     httpPath,
	"httpPathsAdditionalBindings":  httpPathsAdditionalBindings,
	"httpBody":                     httpBody,
	"shortType":                    shortType,
	"urlHasVarsFromMessage":        urlHasVarsFromMessage,
	"lowerGoNormalize":             lowerGoNormalize,
	"goNormalize":                  goNormalize,
	"leadingComment":               leadingComment,
	"trailingComment":              trailingComment,
	"leadingDetachedComments":      leadingDetachedComments,
	"stringFieldExtension":         stringFieldExtension,
	"stringMethodOptionsExtension": stringMethodOptionsExtension,
	"boolMethodOptionsExtension":   boolMethodOptionsExtension,
	"boolFieldExtension":           boolFieldExtension,
	"isFieldMap":                   isFieldMap,
	"fieldMapKeyType":              fieldMapKeyType,
	"fieldMapValueType":            fieldMapValueType,
	"replaceDict":                  replaceDict,
	"getStringSlice":               getStringSlice,
	"getInt":                       getInt,
	"getStringMapString":           getStringMapString,
	"getStringMapInterface":        getStringMapInterface,
	"get":                          get,
	"allSettings":                  allSettings,
	"allKeys":                      allKeys,
	"getTime":                      getTime,
	"getString":                    getString,
}

var getString = func(key string) string {
	return viper.GetString(key)
}

var getStringSlice = func(s string) []string {
	return viper.GetStringSlice(s)
}

var getInt = func(key string) int {
	return viper.GetInt(key)
}

var getStringMapString = func(key string) map[string]string {
	return viper.GetStringMapString(key)
}

var get = func(key string) interface{} {
	return viper.Get(key)
}

var getStringMapInterface = func(key string) interface{} {
	return viper.GetStringMap(key)
}

var allSettings = func() map[string]interface{} {
	return viper.AllSettings()
}

var allKeys = func() []string {
	return viper.AllKeys()
}

var getTime = func(key string) time.Time {
	return viper.GetTime(key)
}
