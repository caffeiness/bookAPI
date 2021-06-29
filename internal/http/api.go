package http

import (
	"bookapi/internal/http/gen"

	"bookapi/internal/http/usecase"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type Api struct {
	books *usecase.BookUsecase
}

func (a Api) FindBooks(ctx echo.Context, params gen.FindBooksParams) error {
	return a.books.FindBooks(ctx, params)
}

func (a Api) AddBook(ctx echo.Context) error {
	return a.books.AddBook(ctx)
}

func (a Api) DeleteBook(ctx echo.Context, id gen.ID) error {
	return a.books.DeleteBook(ctx, id)
}

func (a Api) FindBookById(ctx echo.Context, id gen.ID) error {
	return a.books.FindBookById(ctx, id)
}

func NewApi(db *gorm.DB) *Api {
	return &Api{
		books: usecase.NewBook(),
	}
}

var _ gen.ServerInterface = (*Api)(nil)
