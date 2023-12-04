package transactions

import "time"

type CampaignTransactionFormatter struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionFormatter struct {
	ID        string            `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{
		ID:        transaction.ID,
		Code:      transaction.Code,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	transactionFormatter := make([]CampaignTransactionFormatter, len(transactions))
	for i, transaction := range transactions {
		transactionFormatter[i] = FormatCampaignTransaction(transaction)
	}

	return transactionFormatter

}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign: CampaignFormatter{
			Name:     transaction.Campaign.Name,
			ImageURL: "",
		},
	}

	if len(transaction.Campaign.CampaignImages) > 0 {
		formatter.Campaign.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	return formatter
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	transactionFormatter := make([]UserTransactionFormatter, len(transactions))
	for i, transaction := range transactions {
		transactionFormatter[i] = FormatUserTransaction(transaction)
	}

	return transactionFormatter
}
