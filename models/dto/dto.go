package dto

type UserRegisterModel struct {
	NickName  string
	Email     string
	Password  string
	PromoCode *string
	ReffID    *string
}

type UserLoginModel struct {
	NickName *string
	Email    *string
	Password string
}

type SessionCheck struct {
	Url string
}

type Check_isPublic struct {
	IsPublic bool
	User     interface{}
}

type StringRequest struct {
	Value string
}

type RequestError struct {
	Message    string
	StatusCode uint
	Err        error
}
