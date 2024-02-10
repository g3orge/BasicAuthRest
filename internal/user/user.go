package user

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type srvUser struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }
