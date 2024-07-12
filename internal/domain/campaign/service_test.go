package campaign_test

import (
	"emailn/internal/contracts"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internal-errors"
	internalmock "emailn/internal/test/internal_mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	service     = campaign.ServiceImpl{}
	newCampaign = contracts.NewCampaignDto{
		Name:     "Test Y",
		Content:  "Content",
		Emails:   []string{"test1@test.com"},
		CreateBy: "test@gmail.com",
	}
)

func TestCreateCampaign(t *testing.T) {
	assert := assert.New(t)

	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(nil)

	service.Repository = repositoryMock
	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)
}

func TestCreateCampaignValidateDomainError(t *testing.T) {
	assert := assert.New(t)

	_, err := service.Create(contracts.NewCampaignDto{})

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func TestCreateSaveCampaign(t *testing.T) {
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}

		return true
	})).Return(nil)

	service.Repository = repositoryMock
	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func TestCreateValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)

	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("ErroSaveDatabase"))

	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func TestGetByIdReturnCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreateBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.MatchedBy(func(id string) bool {
		return id == campaign.ID
	})).Return(campaign, nil)

	service.Repository = repositoryMock

	campaignReturned, _ := service.GetById(campaign.ID)

	assert.Equal(campaign.ID, campaignReturned.ID)
	assert.Equal(campaign.Name, campaignReturned.Name)
	assert.Equal(campaign.Content, campaignReturned.Content)
	assert.Equal(campaign.Status, campaignReturned.Status)
	assert.Equal(campaign.CreatedBy, campaignReturned.CreatedBy)
}

func TestGetByIdReturnSomethingWrongExists(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreateBy)

	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(nil, errors.New("Error Something"))

	service.Repository = repositoryMock

	_, err := service.GetById(campaign.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func TestDeleteReturnRecordNotFoundWhenCampaignDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	campaignIdInvalid := "invalid"

	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	service.Repository = repositoryMock

	err := service.Delete(campaignIdInvalid)

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func TestDeleteReturnStatusInvalidWhenCampaignHasStatusNotEqualsPending(t *testing.T) {
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}

	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaign, nil)

	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Equal("Campaign status invalid", err.Error())
}

func TestDeleteReturnInternalErrorWhenDeleteHasProblem(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("Test 1", "Body !!", []string{"teste@teste.com.br"}, newCampaign.CreateBy)

	repositoryMock := new(internalmock.CampaignRepositoryMock)

	repositoryMock.On("GetById", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(errors.New("error to delete campaign"))

	service.Repository = repositoryMock

	err := service.Delete(campaignFound.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func TestDeleteReturnNilWhenDeleteHasSuccess(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("Test 1", "Body !!", []string{"teste@teste.com.br"}, newCampaign.CreateBy)

	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(nil)

	service.Repository = repositoryMock

	err := service.Delete(campaignFound.ID)

	assert.Nil(err)
}
