package http

import (
	"bookapi/internal/http/gen"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type Api struct {
}

func (a Api) FindBooks(ctx echo.Context, params gen.FindBooksParams) error {
	panic("implement me")
}

func (a Api) AddBook(ctx echo.Context) error {
	panic("implement me")
}

func (a Api) DeleteBook(ctx echo.Context, id gen.ID) error {
	panic("implement me")
}

func (a Api) FindBookById(ctx echo.Context, id gen.ID) error {
	panic("implement me")
}

func NewApi(db *gorm.DB) *Api {
	return &Api{}
}

var _ gen.ServerInterface = (*Api)(nil)
