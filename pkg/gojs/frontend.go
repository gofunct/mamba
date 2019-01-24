package gojs

import (
	"github.com/gopherjs/gopherjs/js"
)

type FrontEnd struct{}

func (f *FrontEnd) Write(msg string) {
	js.Global.Get("document").Call("write", msg)
}

func (f *FrontEnd) Alert(msg string) {
	js.Global.Call("alert", msg)
}

func (f *FrontEnd) OnClick(button string, e func()) {
	js.Global.Get(button).Call("addEventListener", "click", func() {
		go func() {
			e()
		}()
	})
}
