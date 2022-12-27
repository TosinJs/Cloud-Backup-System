package userEntity

type UserSignUpReq struct {
	UserId   string `json:"user_id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Status   string `json:"-"`
}

type UserSignUpRes struct {
	Token string `json:"token"`
}

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginRes struct {
	Token string `json:"token"`
}

type UserLoginRepoRes struct {
	Password string `json:"-"`
	Status   string `json:"-"`
}
