package utils

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

func RenderComponentHandler(component templ.Component, handlers ...fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		for _, h := range handlers {
			if err := h(c); err != nil {
				return err
			}
		}
		return RenderComponent(c, component)
	}
}
func RenderComponent(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
