package user

import (
	"github.com/google/uuid"
	"time"
)

// User Model
type User struct {
	ID             uuid.UUID
	Name           string
	Occupation     string
	Email          string
	Role           string
	PasswordHash   string
	Token          string
	AvatarFileName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}
