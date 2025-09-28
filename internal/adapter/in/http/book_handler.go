package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"tabeo.org/challenge/pkg/logger"
)

var validate = validator.New()

type BookDefaultHandler struct {
	log logger.AppLogger
}

func NewBookDefaultHandler(log logger.AppLogger) BookHandler {
	return &BookDefaultHandler{log: log}
}

func (b *BookDefaultHandler) CreateBook(c fiber.Ctx) (*BookResponse, error) {
	var req BookRequest
	if err := c.Bind().Body(&req); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := validate.Struct(&req); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	book := req.ToEntity()
	print(book)

	return nil, c.SendStatus(fiber.StatusNotImplemented)
}
