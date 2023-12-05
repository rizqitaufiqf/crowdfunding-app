package transactions

import "crowdfunding/user"

type GetCampaignTransactionDTO struct {
	ID   string `uri:"id" binding:"required"`
	User user.User
}

type CreateTransactionDTO struct {
	Amount     int    `json:"amount" binding:"required"`
	CampaignID string `json:"campaign_id" binding:"required"`
	User       user.User
}
