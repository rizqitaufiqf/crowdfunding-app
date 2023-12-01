package transactions

import (
	"crowdfunding/user"
	"time"
)

type Transaction struct {
	ID         string
	CampaignID string
	UserID     string
	User       user.User
	Amount     int
	Status     string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
