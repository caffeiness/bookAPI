package http

import (
	"bookapi/internal/database/repository"
	"bookapi/internal/http/gen"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	om "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/go-playground/validator/v10"
)

type Book struct {
	Tag   string `json:"tag"`
	Name  string `json:"name" validate:"required"`
	Price int    `json:"Price" validate:"max=999999999999"`
}

type BookPath struct {
	Tag   *string `json:"tag"`
	Name  *string `json:"name" validate:"required"`
	Price *int    `json:"Price"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func Run() {
	e := echo.New()

	// validator
	spec, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}
	e.Use(om.OapiRequestValidator(spec))

	//mysql connection
	//TODO 設定ファイルの利用と、database共通処理を作る
	dsn := "user:pass@tcp(127.0.0.1:3309)/bookAPI?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	//TODO 外からauto-migration対象を指定できる仕組みを作る
	if err := db.AutoMigrate(&repository.Book{}); err != nil {
		panic(err.Error())
	}

	e.POST("/book", func(context echo.Context) error {
		book := new(Book)
		_ = context.Bind(book)
		if err := context.Validate(book); err != nil {
			return context.String(http.StatusBadRequest, err.Error())
		}
		now := time.Now()
		db.Create(&repository.Book{
			Tag:       book.Tag,
			Name:      book.Name,
			Price:     book.Price,
			CreatedAt: now,
			UpdatedAt: now,
		})
		return context.JSON(http.StatusCreated, book)
	})

	e.GET("/book/:code", func(context echo.Context) error {
		code := context.Param("code")
		m := new(repository.Book)
		if tx := db.First(m, "code=?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}

		book := &Book{
			Tag:   m.Tag,
			Name:  m.Name,
			Price: m.Price,
		}
		return context.JSON(http.StatusOK, book)
	})

	e.PUT("/book/:code", func(context echo.Context) error {
		code := context.Param("code")
		book := new(Book)
		_ = context.Bind(book)
		if err := context.Validate(book); err != nil {
			return context.String(http.StatusBadRequest, err.Error())
		}

		m := new(repository.Book)
		if tx := db.First(m, "code = ?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}
		//アプデ
		now := time.Now()
		db.Model(m).
			Where("code = ?", code).
			Updates(repository.Book{
				Name:      book.Name,
				UpdatedAt: now,
			})
		return context.JSON(http.StatusOK, book)
	})

	e.PATCH("/book/:code", func(context echo.Context) error {
		code := context.Param("code")
		book := new(BookPath)
		_ = context.Bind(book)
		if err := context.Validate(book); err != nil {
			return context.String(http.StatusBadRequest, err.Error())
		}

		m := new(repository.Book)
		if tx := db.First(m, "code = ?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}

		tx := db.Model(m).Where("code = ?", code)
		if book.Price != nil {
			m.Price = *book.Price
		}
		if book.Name != nil {
			m.Name = *book.Name
		}
		tx.Updates(*m)
		return context.JSON(http.StatusOK, &Book{
			Tag:   m.Tag,
			Name:  m.Name,
			Price: m.Price,
		})

	})

	e.DELETE("/book/:code", func(context echo.Context) error {
		code := context.Param("code")
		m := new(repository.Book)
		if tx := db.First(m, "code = ?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}
		db.Delete(m, "code = ?", code)

		return context.String(http.StatusNotFound, "")

	})

	//generateしたhandlerの実装
	gen.RegisterHandlers(e, NewApi(db))
	e.Logger.Fatal(e.Start(":1232"))
}
