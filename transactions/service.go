package transactions

import (
	"crowdfunding/campaign"
	"crowdfunding/helper"
	"crowdfunding/payment"
	"errors"
	"github.com/google/uuid"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionDTO) ([]Transaction, error)
	GetTransactionsByUserID(userID string) ([]Transaction, error)
	CreateTransaction(input CreateTransactionDTO) (Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
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

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}
	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
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
