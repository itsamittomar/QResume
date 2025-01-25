package contracts

type Register struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}


type UserDetails struct{
	Name string `json:"name"`
	Phone string `json:"phone"`
	Linkedin string `json:"linkedin"`
	Github string `json:"github"`
	Leetcode string `json:"leetcode"`
	GeeksForGeeks string `json:"geeksforgeeks"`
	Scaler string `json:"scaler"`
} 