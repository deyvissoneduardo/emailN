package campaign

import (
	"emailn/internal/contracts"
	internalerrors "emailn/internal/internal-errors"
	"errors"
)

type Service interface {
	Create(newCampaign contracts.NewCampaignDto) (string, error)
	GetById(id string) (*contracts.CampaignResponse, error)
	Cancel(id string) error
	Delete(id string) error
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

func (s *ServiceImpl) GetById(id string) (*contracts.CampaignResponse, error) {
	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	return &contracts.CampaignResponse{
		ID:      campaign.ID,
		Name:    campaign.Name,
		Content: campaign.Content,
		Status:  campaign.Status,
	}, nil
}

func (s *ServiceImpl) Cancel(id string) error {

	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return internalerrors.ErrInternal
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	campaign.Cancel()
	err = s.Repository.Update(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}

	return nil
}

func (s *ServiceImpl) Delete(id string) error {

	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return internalerrors.ErrInternal
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	campaign.Delete()
	err = s.Repository.Delete(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}

	return nil
}
