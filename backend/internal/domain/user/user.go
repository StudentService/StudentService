package user

type Role string

const (
	RoleStudent Role = "student"
	// ... другие роли
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string // не json
	Role         Role   `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	// ... другие поля
}

// Value object или методы, если нужно
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}
