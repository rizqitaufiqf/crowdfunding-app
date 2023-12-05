package campaign

import (
	"crowdfunding/user"
	"time"
)

type Campaign struct {
	ID               string
	UserID           string
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time // make this type to a pointer which mean it can be nil value
	CampaignImages   []CampaignImage
	User             user.User
}

type CampaignImage struct {
	ID         string
	CampaignID string
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time // make this type to a pointer which mean it can be nil value
}
