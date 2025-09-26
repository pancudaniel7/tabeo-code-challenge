package apperr

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func HttpHandleError(ctx fiber.Ctx, err error) error {
	var nf *NotFoundErr
	var ae *AlreadyExistsErr
	var ia *InvalidArgErr
	var ve validator.ValidationErrors

	switch {
	case errors.As(err, &nf):
		return ctx.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": nf.Error()})

	case errors.As(err, &ae):
		return ctx.Status(fiber.StatusConflict).
			JSON(fiber.Map{"error": ae.Error()})

	case errors.As(err, &ia):
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": ia.Error()})

	case errors.As(err, &ve):
		errorsMap := make(map[string]string)
		for _, fieldErr := range ve {
			errorsMap[fieldErr.Field()] = fieldErr.Tag()
		}
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"errors": errorsMap})

	default:
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "internal error"})
	}
}
