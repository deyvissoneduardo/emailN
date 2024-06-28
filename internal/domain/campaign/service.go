package campaign

import (
	"emailn/internal/contracts"
	internalerrors "emailn/internal/internal-errors"
)

type Service interface {
	Create(newCampaign contracts.NewCampaignDto) (string, error)
}

type ServiceImpl struct {
	Repository Repository
}

func (s *ServiceImpl) Create(newCampaign contracts.NewCampaignDto) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return "", err
	}

	err = s.Repository.Create(campaign)

	if err != nil {
		return "", internalerrors.ErrInternal
	}

	return campaign.ID, nil
}
