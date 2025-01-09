package ui

import (
	"fmt"

	"github.com/a-h/templ"
)

type AttrModifier func(*templ.Attributes)

func CreateAttrs(baseClass string, opts ...AttrModifier) templ.Attributes {
	attrs := templ.Attributes{
		"class": baseClass,
	}
	for _, o := range opts {
		o(&attrs)
	}
	return attrs
}

func Merge(a, b string) string {
	return fmt.Sprintf("%s %s", a, b)
}

func Attr(name, value string) AttrModifier {
	return func(attrs *templ.Attributes) {
		attr := *attrs
		value := attr[name].(string) + " " + value
		attr[name] = value
	}
}

func OptAttr(name, value string) AttrModifier {
	return func(attrs *templ.Attributes) {
		if value == "" {
			return
		}
		attr := *attrs
		value := attr[name].(string) + " " + value
		attr[name] = value
	}
}

func Class(class string) AttrModifier {
	return Attr("class", class)
}
func Name(name string) AttrModifier {
	return Attr("name", name)
}
