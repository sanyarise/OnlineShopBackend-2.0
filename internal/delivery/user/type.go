package user

type Credentials struct {
	Email    string `json:"email" binding:"required" maxLength:"250" example:"example@gmail.com" format:"email"`
	Password string `json:"password" binding:"required" minLenght:"5" maxLenght:"16"`
}

type UserId struct {
	Value string `json:"id" uri:"userID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type OutUser struct {
	Id        string `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Firstname string `json:"first_name" example:"Иван"`
	Lastname  string `json:"last_name" example:"Иванов"`
	Email     string `json:"email" binding:"required, max=250,email" maxLength:"250" example:"example@gmail.com" format:"email" swaggertype:"string"`
	Address   `json:"address"`
	Rights    `json:"rights"`
}

type Address struct {
	Zipcode string `json:"zipcode,omitempty" example:"000000"`
	Country string `json:"country,omitempty" example:"Russia"`
	City    string `json:"city,omitempty" example:"Moscow"`
	Street  string `json:"street,omitempty" example:"Pushkina"`
}
type Rights struct {
	ID    string   `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Name  string   `json:"rights_name" binding:"required" example:"customer"`
	Rules []string `json:"rules,omitempty"`
}

type ChangeRights struct {
	UserId string
	RightsName string
}
