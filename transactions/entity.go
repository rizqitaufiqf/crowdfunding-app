package transactions

import (
	"crowdfunding/campaign"
	"crowdfunding/user"
	"time"
)

type Transaction struct {
	ID                    string
	CampaignID            string
	UserID                string
	User                  user.User
	Amount                int
	Status                string
	Code                  string
	PaymentURL            string
	MidtransTransactionID string
	Campaign              campaign.Campaign
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}
