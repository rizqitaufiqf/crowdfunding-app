package transactions

import (
	"crowdfunding/campaign"
	"crowdfunding/helper"
	"crowdfunding/user"
	"errors"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"os"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionDTO) ([]Transaction, error)
	GetTransactionsByUserID(userID string) ([]Transaction, error)
	CreateTransaction(input CreateTransactionDTO) (Transaction, error)
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input TransactionNotificationDTO) error
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionDTO) ([]Transaction, error) {
	camp, err := s.campaignRepository.FindByCampaignID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if camp.UserID != input.User.ID {
		return []Transaction{}, errors.New("invalid user")
	}

	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *service) GetTransactionsByUserID(userID string) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionDTO) (Transaction, error) {
	transaction := Transaction{
		ID:         uuid.New().String(),
		CampaignID: input.CampaignID,
		Amount:     input.Amount,
		UserID:     input.User.ID,
		Status:     "pending",
		Code:       helper.GenerateTransactionCode(),
	}

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}
	paymentURL, err := s.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL
	updatedTransaction, err := s.repository.Update(newTransaction)
	if err != nil {
		return updatedTransaction, err
	}

	return updatedTransaction, nil
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	var sn = snap.Client{}
	sn.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.ID,
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapResp, err := sn.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return snapResp.RedirectURL, nil
}

func (s *service) ProcessPayment(input TransactionNotificationDTO) error {
	transactionID := input.OrderID

	transaction, err := s.repository.GetByTransactionID(transactionID)
	if err != nil {
		return err
	}

	transaction.MidtransTransactionID = input.TransactionID
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "canceled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	camp, err := s.campaignRepository.FindByCampaignID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		camp.BackerCount = camp.BackerCount + 1
		camp.CurrentAmount = camp.CurrentAmount + updatedTransaction.Amount

		if _, err := s.campaignRepository.UpdateCampaign(camp); err != nil {
			return err
		}
	}

	return nil
}
