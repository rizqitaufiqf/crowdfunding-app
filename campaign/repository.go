package campaign

import (
	"errors"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Campaign, error)
	FindAllByUserID(UserId string) ([]Campaign, error)
	FindByCampaignID(campaignID string) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	UpdateCampaign(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	// Preload "CampaignImages" association with the condition "campaign_images.is_primary = 1"
	// This ensures that related CampaignImages are eagerly loaded along with the Campaigns during the Find operation.
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindAllByUserID(userID string) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByCampaignID(campaignID string) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", campaignID).Find(&campaign).Error
	if err != nil {
		return campaign, err
	}

	if campaign.ID == "" || campaign.ID == "00000000-0000-0000-0000-000000000000" {
		return campaign, errors.New("campaign not found")
	}

	return campaign, nil
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) UpdateCampaign(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
