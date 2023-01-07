package cart

type Cart struct {
	Id     string     `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	UserId string     `json:"user_id,omitempty" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Items  []CartItem `json:"items" binding:"min=0" minimum:"0"`
}

type ShortCart struct {
	CartId string `json:"cart_id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	ItemId string `json:"item_id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type CartId struct {
	Value string `json:"id" uri:"cartID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type CartItem struct {
	Id    string `json:"item_id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Title string `json:"title" binding:"required" example:"Пылесос"`
	Price int32  `json:"price" example:"1990" default:"10" binding:"required" minimum:"0"`
	Image string `json:"image,omitempty"`
}
