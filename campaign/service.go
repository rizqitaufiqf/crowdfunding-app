package campaign

import (
	"crowdfunding/helper"
	"errors"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"strings"
)

type Service interface {
	GetCampaigns(UserID string) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDTO) (Campaign, error)
	CreateCampaign(input CreateCampaignDTO) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDTO, inputData CreateCampaignDTO) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(UserID string) ([]Campaign, error) {
	if UserID != "" && UserID != "00000000-0000-0000-0000-000000000000" {
		campaigns, err := s.repository.FindAllByUserID(UserID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *service) GetCampaignByID(input GetCampaignDTO) (Campaign, error) {
	campaign, err := s.repository.FindByCampaignID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignDTO) (Campaign, error) {
	campaign := Campaign{
		ID:               uuid.New().String(),
		UserID:           input.User.ID,
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.ShortDescription,
		GoalAmount:       input.GoalAmount,
		Perks:            helper.SanitizePerksSplitString(input.Perks),
		Slug:             slug.Make(input.Name + "-" + strings.Split(uuid.New().String(), "-")[4]),
	}

	campaign, err := s.repository.Save(campaign)
	if err != nil {
		return campaign, err
	}

	return campaign, err
}

func (s *service) UpdateCampaign(inputID GetCampaignDTO, inputData CreateCampaignDTO) (Campaign, error) {
	campaign, err := s.repository.FindByCampaignID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("invalid user")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.GoalAmount = inputData.GoalAmount
	campaign.Perks = inputData.Perks

	updatedCampaign, err := s.repository.UpdateCampaign(campaign)
	if err != nil {
		return updatedCampaign, nil
	}

	return updatedCampaign, nil
}
