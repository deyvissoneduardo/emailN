package campaign

import (
	"emailn/internal/contracts"
	internalerrors "emailn/internal/internal-errors"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	service     = Service{}
	newCampaign = contracts.NewCampaignDto{
		Name:    "Test Y",
		Content: "Content",
		Emails:  []string{"test1@test.com"},
	}
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Create(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func TestCreateCampaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(repositoryMock)
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
	repositoryMock := new(repositoryMock)
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *Campaign) bool {
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

	repositoryMock := new(repositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("ErroSaveDatabase"))

	service.Repository = repositoryMock
	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}
