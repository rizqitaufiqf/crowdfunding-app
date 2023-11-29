package campaign

type Service interface {
	GetCampaigns(UserID string) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
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

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByCampaignID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
