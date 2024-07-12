package campaign

import (
	"emailn/internal/contracts"
	internalerrors "emailn/internal/internal-errors"
	"errors"
)

type Service interface {
	Create(newCampaign contracts.NewCampaignDto) (string, error)
	GetById(id string) (*contracts.CampaignResponse, error)
	Delete(id string) error
}

type ServiceImpl struct {
	Repository Repository
}

func (s *ServiceImpl) Create(newCampaign contracts.NewCampaignDto) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreateBy)

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
		return nil, internalerrors.ProcessErroToReturn(err)
	}

	return &contracts.CampaignResponse{
		ID:                   campaign.ID,
		Name:                 campaign.Name,
		Content:              campaign.Content,
		Status:               campaign.Status,
		AmountOfEmailsToSend: len(campaign.Contacts),
		CreatedBy:            campaign.CreatedBy,
	}, nil
}

func (s *ServiceImpl) Delete(id string) error {

	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return internalerrors.ProcessErroToReturn(err)
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	campaign.Delete()
	err = s.Repository.Delete(campaign)
	if err != nil {
		return internalerrors.ProcessErroToReturn(err)
	}

	return nil
}
