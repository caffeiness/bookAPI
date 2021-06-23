package http

import (
	"bookapi/internal/database/repository"
	"bookapi/internal/http/gen"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	om "github.com/deepmap/oapi-codegen/pkg/middleware"
)

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
	dsn := "user:pass@tcp(127.0.0.1:3309)/bookApi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	//TODO 外からauto-migration対象を指定できる仕組みを作る
	if err := db.AutoMigrate(&repository.Book{}); err != nil {
		panic(err.Error())
	}

	//generateしたhandlerの実装
	gen.RegisterHandlers(e, NewApi(db))

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.Logger.Fatal(e.Start(":1232"))
}
