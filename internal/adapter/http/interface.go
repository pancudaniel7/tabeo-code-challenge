package http

import "github.com/gofiber/fiber/v3"

type BookHandler interface {
	CreateBook(ctx fiber.Ctx) (*BookResponse, error)
}
