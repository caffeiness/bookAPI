package usecase

import (
	"bookapi/internal/http/gen"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/labstack/echo/v4"
)

type BookUsecase struct {
	Books  map[int64]gen.BookResponse
	NextId int64
	Lock   sync.Mutex
}

func NewBook() *BookUsecase {
	return &BookUsecase{
		Books:  make(map[int64]gen.BookResponse),
		NextId: 1000,
	}
}

func (p *BookUsecase) FindBooks(ctx echo.Context, params gen.FindBooksParams) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	var result []gen.BookResponse

	// 順番通りに並べる(デフォルトAsc)
	isAsc := true
	if params.Order != nil && *params.Order == "desc" {
		isAsc = false
	}
	fmt.Println(isAsc)
	keys := make([]int64, 0, len(p.Books))
	for k := range p.Books {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		if isAsc {
			return keys[gen.ID(i)] < keys[gen.ID(j)]
		}
		return keys[gen.ID(j)] < keys[gen.ID(i)]
	})

	fmt.Println(keys)

	for _, key := range keys {
		pet := p.Books[key]
		if params.Tags != nil {
			// If we have tags,  filter pets by tag
			for _, t := range *params.Tags {
				if pet.Tag != nil && (*pet.Tag == t) {
					result = append(result, pet)
				}
			}
		} else {
			// Add all pets if we're not filtering
			result = append(result, pet)
		}

		if params.Limit != nil {
			l := int(*params.Limit)
			if len(result) >= l {
				// We're at the limit
				break
			}
		}
	}
	return ctx.JSON(http.StatusOK, result)
}

func (p *BookUsecase) AddBook(ctx echo.Context) error {
	req := new(gen.Book)
	err := ctx.Bind(&req)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for NewBook")
	}
	// We now have a np, let's add it to our "database".

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	// We handle pets, not NewBooks, which have an additional ID field
	res := gen.BookResponse{
		Book: gen.Book{
			Id:   &p.NextId,
			Name: req.Name,
			Tag:  req.Tag,
		},
	}
	p.Books[p.NextId] = res
	p.NextId = p.NextId + 1

	// Now, we have to return the NewBook
	err = ctx.JSON(http.StatusCreated, res)
	if err != nil {
		// Something really bad happened, tell Echo that our handler failed
		return err
	}

	// Return no error. This refers to the handler. Even if we return an HTTP
	// error, but everything else is working properly, tell Echo that we serviced
	// the error. We should only return errors from Echo handlers if the actual
	// servicing of the error on the infrastructure level failed. Returning an
	// HTTP/400 or HTTP/500 from here means Echo/HTTP are still working, so
	// return nil.
	return nil
}

func (p *BookUsecase) FindBookById(ctx echo.Context, id gen.ID) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	pet, found := p.Books[id.Int64()]
	if !found {
		return sendError(ctx, http.StatusNotFound,
			fmt.Sprintf("Could not find pet with ID %d", id))
	}
	return ctx.JSON(http.StatusOK, pet)
}

func (p *BookUsecase) DeleteBook(ctx echo.Context, id gen.ID) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, found := p.Books[id.Int64()]
	if !found {
		return sendError(ctx, http.StatusNotFound,
			fmt.Sprintf("Could not find pet with ID %d", id))
	}
	delete(p.Books, id.Int64())
	return ctx.NoContent(http.StatusNoContent)
}
