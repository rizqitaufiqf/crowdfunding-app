package campaign

import "time"

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
	DeletedAt        time.Time
	CampaignImages   []CampaignImage
}

type CampaignImage struct {
	ID         string
	CampaignID string
	FileName   string
	IsPrimary  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}
