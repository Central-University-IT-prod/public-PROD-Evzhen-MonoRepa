package entities

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Login    string `json:"login"`
	Password []byte `json:"password"`
	//Profiles []Profile `json:"profiles" gorm:"foreignKey:UserLogin; constraint:OnDelete:CASCADE"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
