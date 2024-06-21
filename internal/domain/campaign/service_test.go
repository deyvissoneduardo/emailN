package campaign

import (
	"emailn/internal/contracts"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	service := Service{}

	newCampaign := contracts.NewCampaignDto{
		Name:    "Test Y",
		Content: "Content",
		Emails:  []string{"test1@test.com"},
	}

	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)
}

func TestCreateSaveCampaign(t *testing.T) {

	newCampaign := contracts.NewCampaignDto{
		Name:    "Test Y",
		Content: "sContent",
		Emails:  []string{"test1@test.com"},
	}

	repositoryMock := new(repositoryMock)
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}

		return true
	})).Return(nil)

	service := Service{Repository: repositoryMock}

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)

}
