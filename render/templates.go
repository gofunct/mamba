package render

import (
	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"
	"io"
	"text/template"
)

var TemplateFuncs = funcMap()

// tmpl executes the given template text on data, writing the result to w.
func Compile(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(TemplateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

func funcMap() template.FuncMap {
	newMap := sprig.GenericFuncMap()
	for k, v := range viper.AllSettings() {
		newMap[k] = v
	}
	return newMap
}

func addTemplateFunc(name string, tmplFunc interface{}) {
	TemplateFuncs[name] = tmplFunc
}

func addTemplateFuncs(tmplFuncs template.FuncMap) {
	for k, v := range tmplFuncs {
		TemplateFuncs[k] = v
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
