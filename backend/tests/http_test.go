package tests

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/go-resty/resty/v2"
	"testing"
)

type Product struct {
	Name      string `json:"name,omitempty"`
	Size      string `json:"size,omitempty"`
	ProductId int    `json:"product_id,omitempty"`
	Count     int    `json:"count,omitempty"`
}

func TestGetHandler(t *testing.T) {
	shouldBeGoods := []Product{
		{
			Name:      "milk",
			Size:      "0.2m",
			ProductId: 1,
			Count:     15},
		{
			Name:      "tables",
			Size:      "1.5m",
			ProductId: 2,
			Count:     20,
		},
	}

	var goods struct {
		AllCount int       `json:"all_count,omitempty"`
		Goods    []Product `json:"goods,omitempty"`
	}

	shouldBeAllCount := 35

	shouldBeGoods = append(shouldBeGoods)

	client := resty.New()

	resp, _ := client.R().Get("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &goods)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, shouldBeGoods, goods.Goods)

	assert.Equal(t, shouldBeAllCount, goods.AllCount)
}

func TestPostHandler(t *testing.T) {
	shouldBeGoods := []Product{
		{
			Name:      "milk",
			Size:      "0.2m",
			ProductId: 1,
			Count:     15},
		{
			Name:      "tables",
			Size:      "1.5m",
			ProductId: 2,
			Count:     20,
		},
	}

	client := resty.New()

	var getData []Product

	resp, _ := client.R().
		SetBody(`[1,2]`).
		Post("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &getData)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, shouldBeGoods, getData)

	var respString string

	resp, _ = client.R().
		SetBody(`[]`).
		Post("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &respString)

	assert.Equal(t, "Массив не должен быть пустым", respString)
	assert.Equal(t, 400, resp.StatusCode())

	resp, _ = client.R().
		SetBody(`[123123123, 123123123]`).
		SetResult(respString).
		Post("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &respString)

	assert.Equal(t, "Таких товаров нет на складе", respString)
	assert.Equal(t, 400, resp.StatusCode())

	resp, _ = client.R().
		SetBody(`asdasdasd`).
		SetResult(respString).
		Post("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &respString)

	assert.Equal(t, "Такой массив не может быть обработан", respString)
	assert.Equal(t, 400, resp.StatusCode())
}

func TestDeleteHandler(t *testing.T) {
	shouldBeGoods := []Product{
		{
			Name:      "milk",
			Size:      "0.2m",
			ProductId: 1,
			Count:     15},
		{
			Name:      "tables",
			Size:      "1.5m",
			ProductId: 2,
			Count:     20,
		},
	}

	client := resty.New()

	var getData []Product

	resp, _ := client.R().
		SetBody(`[1,2]`).
		SetResult(getData).
		Delete("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &getData)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, shouldBeGoods, getData)

	var respString string

	resp, _ = client.R().
		SetBody(`[]`).
		SetResult(respString).
		Post("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &respString)

	assert.Equal(t, "Массив не должен быть пустым", respString)
	assert.Equal(t, 400, resp.StatusCode())

	resp, _ = client.R().
		SetBody(`[123123123, 123123123]`).
		SetResult(respString).
		Post("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &respString)

	assert.Equal(t, "Таких товаров нет на складе", respString)
	assert.Equal(t, 400, resp.StatusCode())

	resp, _ = client.R().
		SetBody(`asdasdasd`).
		SetResult(respString).
		Post("http://localhost:4000/product")

	json.Unmarshal(resp.Body(), &respString)

	assert.Equal(t, "Такой массив не может быть обработан", respString)
	assert.Equal(t, 400, resp.StatusCode())

}
