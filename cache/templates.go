package cache

import (
	"io"
	"text/template"
)

func (c *cache) AddTmplConfig(src, dest, pkg string) {
	c.v.SetDefault("template.dir", src)
	c.v.SetDefault("template.dest", dest)
	c.v.SetDefault("template.debug", true)
	c.v.SetDefault("template.single_pkg", false)
	c.v.SetDefault("template.all", false)
	c.v.SetDefault("template.pkg", pkg)
}

// tmpl executes the given template text on data, writing the result to w.
func (c *cache) Compile(w io.Writer, text string) error {
	t := template.New("")
	t.Funcs(c.fmap)
	template.Must(t.Parse(text))
	return t.Execute(w, c.v.AllSettings())
}

func (c *cache) AddTmplFunc(name string, tmplFunc interface{}) {
	c.fmap[name] = tmplFunc
}

func (c *cache) AddTmplFuncs(tmplFuncs template.FuncMap) {
	for k, v := range tmplFuncs {
		c.fmap[k] = v
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
