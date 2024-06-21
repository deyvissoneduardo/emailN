package campaign

import (
	"emailn/internal/contracts"
)

type Service struct {
	Repository Repository
}

func (s *Service) Create(newCampaign contracts.NewCampaignDto) (string, error) {

	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	s.Repository.Create(campaign)

	return campaign.ID, nil
}
