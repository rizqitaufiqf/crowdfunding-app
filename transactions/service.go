package transactions

import (
	"crowdfunding/campaign"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionDTO) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionDTO) ([]Transaction, error) {
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
