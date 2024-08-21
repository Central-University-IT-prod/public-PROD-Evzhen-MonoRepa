package entities

import "golang.org/x/crypto/bcrypt"

type Admin struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Login    string `json:"login"`
	Password []byte `json:"password"`
	//Contests []Contest `json:"contests" gorm:"foreignKey:AdminID; constraint:OnDelete:CASCADE"`
}

func (user *Admin) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}
func (user *Admin) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
