package transactions

import (
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID string) ([]Transaction, error)
	GetByUserID(userID string) ([]Transaction, error)
	Save(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(campaignID string) ([]Transaction, error) {
	var transactions []Transaction
	if err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("created_at desc").Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetByUserID(userID string) ([]Transaction, error) {
	var transactions []Transaction
	if err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) Save(transaction Transaction) (Transaction, error) {
	if err := r.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	if err := r.db.Save(&transaction).Error; err != nil {
		return transaction, err
	}
	
	return transaction, nil
}
