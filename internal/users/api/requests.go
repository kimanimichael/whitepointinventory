package api

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email_address"`
	Password string `json:"password"`
}
