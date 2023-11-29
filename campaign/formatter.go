package campaign

import "strings"

type CampaignFormatter struct {
	ID               string `json:"id"`
	UserID           string `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignDetailImages struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type CampaignDetailUser struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignDetailFormatter struct {
	ID               string                 `json:"id"`
	UserID           string                 `json:"user_id"`
	Name             string                 `json:"name"`
	ShortDescription string                 `json:"short_description"`
	Description      string                 `json:"description"`
	CampaignUser     CampaignDetailUser     `json:"user"`
	ImageURL         string                 `json:"image_url"`
	GoalAmount       int                    `json:"goal_amount"`
	CurrentAmount    int                    `json:"current_amount"`
	Slug             string                 `json:"slug"`
	Perks            []string               `json:"perks"`
	CampaignImages   []CampaignDetailImages `json:"images"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	formatter := CampaignFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		ImageURL:         "",
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
	}

	if len(campaign.CampaignImages) > 0 {
		formatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return formatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	campaignsFormatter := make([]CampaignFormatter, len(campaigns))
	for i, campaign := range campaigns {
		campaignsFormatter[i] = FormatCampaign(campaign)
	}

	return campaignsFormatter
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	userFormatter := CampaignDetailUser{
		Name:     campaign.User.Name,
		ImageURL: campaign.User.AvatarFileName,
	}

	images := make([]CampaignDetailImages, len(campaign.CampaignImages))
	for i, image := range campaign.CampaignImages {
		images[i] = CampaignDetailImages{
			ImageURL:  image.FileName,
			IsPrimary: image.IsPrimary == 1,
		}
	}

	campaignDetailFormatter := CampaignDetailFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		CampaignUser:     userFormatter,
		ImageURL:         "",
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		Perks:            strings.Split(campaign.Perks, ";"),
		CampaignImages:   images,
	}

	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignDetailFormatter
}
