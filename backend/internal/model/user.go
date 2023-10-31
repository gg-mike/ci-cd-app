package model

// TODO: rbac
type UserCore struct {
	Username string `json:"username" gorm:"uniqueIndex"`
}

type User struct {
	UserCore
	Common
}

type UserCreate struct {
	UserCore
}

type UserShort struct {
	UserCore
	Common
}