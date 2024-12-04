package htmx

import (
	"github.com/FrancescoLuzzi/AQuickQuestion/app/app_context"
	"github.com/gofiber/fiber/v3"
)

func New() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		ctx.Locals(app_context.LayoutCtxKey, ctx.Get("hx-request", "false") == "true")
		return nil
	}
}
