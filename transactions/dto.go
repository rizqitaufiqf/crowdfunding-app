package transactions

import "crowdfunding/user"

type GetCampaignTransactionDTO struct {
	ID   string `uri:"id" binding:"required"`
	User user.User
}
