package usecase

import (
	"bookapi/internal/http/gen"

	"github.com/labstack/echo/v4"
)

func sendError(ctx echo.Context, code int, message string) error {
	petErr := gen.Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}
