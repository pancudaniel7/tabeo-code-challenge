package apperr

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func HttpHandleError(ctx fiber.Ctx, err error) error {
	var ve validator.ValidationErrors

	switch {
	case IsNotFound(err):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": messageOf(err)})
	case IsExists(err):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": messageOf(err)})
	case IsInvalidArgument(err):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": messageOf(err)})
	case IsInvalid(err):
		if errors.As(err, &ve) {
			errorsMap := make(map[string]string)
			for _, fe := range ve {
				errorsMap[fe.Field()] = fe.Tag()
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errorsMap})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": messageOf(err)})
	case errors.As(err, &ve):
		errorsMap := make(map[string]string)
		for _, fe := range ve {
			errorsMap[fe.Field()] = fe.Tag()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errorsMap})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}
}

func messageOf(err error) string {
	var ce *Error
	if errors.As(err, &ce) {
		return ce.Message()
	}
	return err.Error()
}
