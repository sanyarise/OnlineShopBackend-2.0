package category

type ShortCategory struct {
	Name        string `json:"name" binding:"required" example:"Электротехника"`
	Description string `json:"description" binding:"required" example:"Электротехнические товары для дома"`
}

type CategoryId struct {
	Value string `json:"id" uri:"categoryID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type Category struct {
	Id          string `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Name        string `json:"name" binding:"required" example:"Электротехника"`
	Description string `json:"description" binding:"required" example:"Электротехнические товары для дома"`
	Image       string `json:"image,omitempty"`
}

type CategoriesList struct {
	List []Category `json:"categories" binding:"min=0" minimum:"0"`
}