package handlers

import "github.com/gofiber/fiber/v2"

type response struct {
	Msg  string       `json:"message"`
	Data *interface{} `json:"data,omitempty"`
}

func reply(ctx *fiber.Ctx, sc int, msg string, data interface{}) (err error) {
	ctx.Status(sc)
	if data == nil {
		err = ctx.JSON(response{Msg: msg})
	} else {
		err = ctx.JSON(response{Msg: msg, Data: &data})
	}
	return
}
