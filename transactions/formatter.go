package transactions

import "time"

type CampaignTransactionFormatter struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
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
