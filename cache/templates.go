package cache

import (
	"github.com/Masterminds/sprig"
	"io"
	"text/template"
)

func init() {
	TemplateFuncs = Cache.funcMap()
}

var TemplateFuncs template.FuncMap

// tmpl executes the given template text on data, writing the result to w.
func (c *cache) Compile(w io.Writer, text string) error {
	t := template.New("")
	t.Funcs(TemplateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, c.v.AllSettings())
}

func (c *cache) funcMap() template.FuncMap {
	newMap := sprig.GenericFuncMap()
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
