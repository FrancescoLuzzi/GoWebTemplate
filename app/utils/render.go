package utils

import (
	"net/http"

	"github.com/a-h/templ"
)

func RenderComponentHandler(component templ.Component) http.Handler {
	return templ.Handler(component, templ.WithStreaming())
}
