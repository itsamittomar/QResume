package contracts

type Register struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserDetails struct {
	Email         string `json:"email" binding:"required,email"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Linkedin      string `json:"linkedin"`
	Github        string `json:"github"`
	Leetcode      string `json:"leetcode"`
	GeeksForGeeks string `json:"geeksforgeeks"`
	Scaler        string `json:"scaler"`
	Password string `json:"password" binding:"required,min=8"`
}
