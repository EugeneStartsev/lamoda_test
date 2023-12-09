package main

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

type httpServer struct {
	db     *goqu.Database
	router *gin.Engine
}

type Product struct {
	Name      string `json:"name,omitempty" db:"product_name"`
	Size      string `json:"size,omitempty" db:"size"`
	ProductId int    `json:"product_id,omitempty" db:"unique_code"`
	Count     int    `json:"count,omitempty" db:"count"`
}

func newHttpServer(db *goqu.Database) *httpServer {
	s := httpServer{
		db:     db,
		router: gin.Default(),
	}

	s.router.GET("/product", s.handleGetProducts)
	s.router.POST("/product", s.handlePostProducts)
	s.router.DELETE("/product", s.handleDeleteProducts)

	return &s
}

func (s *httpServer) run(listenAddr string) error {
	return s.router.Run(listenAddr)
}

func (s *httpServer) handleGetProducts(c *gin.Context) {
	var query struct {
		WarehouseName string `form:"warehouse"`
	}

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var goods struct {
		AllCount int       `json:"all_count,omitempty" db:"sum"`
		Goods    []Product `json:"goods,omitempty"`
	}

	_, err = s.db.From(goqu.L("products where not exists ?",
		s.db.From("warehouse").Where(goqu.L("products.product_name = warehouse.product_name")))).
		Select(goqu.SUM("count")).ScanVal(&goods.AllCount)

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = s.db.
		From(goqu.L("products where not exists ?",
			s.db.From("warehouse").Where(goqu.L("products.product_name = warehouse.product_name")))).
		ScanStructs(&goods.Goods)

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, goods)
}

func (s *httpServer) handlePostProducts(c *gin.Context) {
	var product Product
	var goods []Product

	query, isExit := bindJson(c)

	var countOfAdd int

	if isExit {
		return
	}

	for id := range query {
		canFindID, err := s.db.From("products").Where(goqu.Ex{"unique_code": query[id]}).ScanStruct(&product)

		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if canFindID {
			var productName string

			canFindName, err := s.db.From("warehouse").Select("product_name").
				Where(goqu.Ex{"product_name": product.Name}).ScanVal(&productName)

			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}

			if !canFindName {
				_, err = s.db.Insert("warehouse").Rows(
					goqu.Record{
						"product_name": product.Name,
						"can_be_use":   true,
					}).Executor().Exec()

				if err != nil {
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				}

				goods = append(goods, product)
			} else {
				countOfAdd++
			}
		}
	}

	if len(goods) == 0 && countOfAdd != 0 {
		c.JSON(http.StatusBadRequest, "Запрашиваемые товары уже зарезервированы или их нет в списке товаров")
	} else if countOfAdd == 0 && len(goods) == 0 {
		c.JSON(http.StatusBadRequest, "Таких товаров нет на складе")
	} else if len(goods) > 0 {
		c.JSON(http.StatusOK, goods)
	}
}

func (s *httpServer) handleDeleteProducts(c *gin.Context) {
	var product Product
	var goods []Product

	query, isExit := bindJson(c)

	if isExit {
		return
	}

	for id := range query {
		canFindID, err := s.db.Select("p.product_name", "p.size", "p.unique_code", "p.count").
			From(goqu.T("products").As("p")).
			Join(goqu.T("warehouse").As("w"), goqu.On(goqu.Ex{"p.product_name": goqu.I("w.product_name")})).
			Where(goqu.Ex{"unique_code": query[id]}).
			ScanStruct(&product)

		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if canFindID {
			_, err = s.db.Delete("warehouse").Where(goqu.Ex{"product_name": product.Name}).Executor().Exec()

			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			goods = append(goods, product)
		}
	}

	if len(goods) == 0 {
		c.JSON(http.StatusBadRequest, "Запрашиваемые товары уже удалены из резерва или не были туда добавлены")
	} else {
		c.JSON(http.StatusOK, goods)
	}

}

func bindJson(c *gin.Context) ([]int, bool) {
	var query []int

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, "Такой массив не может быть обработан")
		return nil, true
	}

	if len(query) == 0 {
		c.JSON(http.StatusBadRequest, "Массив не может быть пустым")
		return nil, true
	} else {
		return query, false
	}
}
