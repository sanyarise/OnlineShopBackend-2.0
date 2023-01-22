package rights

type ShortRights struct {
	Name  string   `json:"name" binding:"required" example:"admin"`
	Rules []string `json:"rules,omitempty"`
}

type OutRights struct {
	Id    string   `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Name  string   `json:"name" binding:"required" example:"admin"`
	Rules []string `json:"rules,omitempty"`
}

type RightsId struct {
	Value string `json:"id" uri:"itemID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}