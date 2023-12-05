package user

import "time"

// User Model
type User struct {
	ID             string
	Name           string
	Occupation     string
	Email          string
	Role           string
	PasswordHash   string
	AvatarFileName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time // make this type to a pointer which mean it can be nil value
}
