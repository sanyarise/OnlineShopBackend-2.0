package item

import "OnlineShopBackend/internal/delivery/category"

type ShortItem struct {
	Title       string `json:"title" binding:"required" example:"Пылесос"`
	Description string `json:"description" binding:"required" example:"Мощность всасывания 1.5 кВт"`
	Category    string `json:"category" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Price       int32  `json:"price" example:"1990" default:"0" binding:"min=0" minimum:"0"`
	Vendor      string `json:"vendor" binding:"required" example:"Витязь"`
}

type ItemId struct {
	Value string `json:"id" uri:"itemID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type Item struct {
	Id          string `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Title       string `json:"title" binding:"required" example:"Пылесос"`
	Description string `json:"description" binding:"required" example:"Мощность всасывания 1.5 кВт"`
	Category    category.Category
	Price       int32    `json:"price" example:"1990" default:"0" binding:"min=0" minimum:"0"`
	Vendor      string   `json:"vendor" binding:"required" example:"Витязь"`
	Images      []string `json:"images,omitempty"`
}

type ItemsList struct {
	List []Item `json:"items" binding:"min=0" minimum:"0"`
}
