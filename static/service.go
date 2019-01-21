package static

import (
	"github.com/shiyanhui/hero"
)

type Static struct{}

func NewStatic() *Static {
	return &Static{}
}

func (s *Static) Generate(source, dest, pkg string) {
	hero.Generate(source, dest, pkg)
}
