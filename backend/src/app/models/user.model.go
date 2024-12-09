package models

import "github.com/cmo7/folly4/src/lib/generics/common"

type UserEntity struct {
	BaseModel `gorm:"embedded"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
}

func (u *UserEntity) GetEntityName() common.EntityName {
	return common.EntityName("User")
}

func (u *UserEntity) GetName() string {
	return u.Username
}
