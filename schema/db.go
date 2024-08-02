package schema

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	Password       string `json:"password"`
}
