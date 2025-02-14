package models

type User struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Email    string  `json:"email" gorm:"unique;not null"`
	Password string  `json:"password" gorm:"not null"`
	Name     string  `json:"name" gorm:"type:varchar(40);not null"`
	QR       []byte  `json:"qr" gorm:"type:bytea"` // Fixed: Changed from blob to bytea
	Details  Details `json:"details" gorm:"embedded;embeddedPrefix:details_"`
}

type Details struct {
	Email         string `json:"email"` // Removed unique and not null constraints
	Phone         string `json:"phone" gorm:"type:varchar(20)"`
	Linkedin      string `json:"linkedin" gorm:"type:varchar(255)"`
	Github        string `json:"github" gorm:"type:varchar(255)"`
	Leetcode      string `json:"leetcode" gorm:"type:varchar(255)"`
	GeeksForGeeks string `json:"geeksforgeeks" gorm:"type:varchar(255)"`
	Scaler        string `json:"scaler" gorm:"type:varchar(255)"`
	QRCodeURL     string `json:"qrCodeUrl" gorm:"type:varchar(255)"`
}
