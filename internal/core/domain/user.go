package domain

type User struct {
	ID       int
	Email    string
	Password string
	Username string
	Role     string
}

type UserReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AddUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
