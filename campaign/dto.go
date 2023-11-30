package campaign

import "crowdfunding/user"

type GetCampaignDTO struct {
	ID string `uri:"id" binding:"required"`
}

type CreateCampaignDTO struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"gt=0,required"`
	Perks            string `json:"perks" binding:"required"`
	User             user.User
}

type CreateCampaignImageDTO struct {
	CampaignID string `form:"campaign_id" binding:"required"`
	IsPrimary  *bool  `form:"is_primary" binding:"required,boolean"`
	User       user.User
}
