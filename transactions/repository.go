package transactions

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID string) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(campaignID string) ([]Transaction, error) {
	var transaction []Transaction
	if err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("created_at desc").Find(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}
