package campaign

type GetCampaignDTO struct {
	ID string `uri:"id" binding:"required"`
}
