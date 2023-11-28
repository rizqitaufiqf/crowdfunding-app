package campaign

type Service interface {
	GetCampaigns(UserID string) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(UserID string) ([]Campaign, error) {
	if UserID != "" && UserID != "00000000-0000-0000-0000-000000000000" {
		campaigns, err := s.repository.FindAllByID(UserID)
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
