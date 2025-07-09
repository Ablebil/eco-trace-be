package response

import "github.com/gofiber/fiber/v2"

func OK(ctx *fiber.Ctx, payload any, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}

	return ctx.Status(fiber.StatusOK).JSON(Res{
		StatusCode: fiber.StatusOK,
		Message:    msg,
		Payload:    payload,
	})
}

func Created(ctx *fiber.Ctx, payload any, message ...string) error {
	msg := "Resource created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	return ctx.Status(fiber.StatusCreated).JSON(Res{
		StatusCode: fiber.StatusCreated,
		Message:    msg,
		Payload:    payload,
	})
}
