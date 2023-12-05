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

type TransactionNotificationDTO struct {
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
