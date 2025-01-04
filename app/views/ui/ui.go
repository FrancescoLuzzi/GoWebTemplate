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

func Class(class string) AttrModifier {
	return func(attrs *templ.Attributes) {
		attr := *attrs
		class := attr["class"].(string) + " " + class
		attr["class"] = class
	}
}
