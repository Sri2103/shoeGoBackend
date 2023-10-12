package userModel

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
	Username  string `json:"username"`
	TokenHash string `json:"tokenHash,omitempty"`
}
