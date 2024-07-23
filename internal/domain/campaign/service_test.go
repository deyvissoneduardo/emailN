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
	newCampaign = contracts.NewCampaignDto{
		Name:     "Test Y",
		Content:  "Body Hi!",
		Emails:   []string{"teste1@test.com"},
		CreateBy: "teste@teste.com.br",
	}
	campaignPending *campaign.Campaign
	campaignStarted *campaign.Campaign
	repositoryMock  *internalmock.CampaignRepositoryMock
	service         = campaign.ServiceImpl{}
)

func setUp() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	campaignPending, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreateBy)
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}
}

func setUpGetByIdRepositoryBy(campaign *campaign.Campaign) {
	repositoryMock.On("GetById", mock.Anything).Return(campaign, nil)
}

func setUpUpdateRepository() {
	repositoryMock.On("Update", mock.Anything).Return(nil)
}

func setUpSendEmailWithSuccess() {
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendEmail = sendMail
}

func TestCreateRequestIsValidIdIsNotNil(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(t, id)
	assert.Nil(t, err)
}

func TestCreateRequestIsNotValidErrInternal(t *testing.T) {
	setUp()

	_, err := service.Create(contracts.NewCampaignDto{})

	assert.False(t, errors.Is(internalerrors.ErrInternal, err))
}

func TestCreateRequestIsValidCallRepository(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}

		return true
	})).Return(nil)

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func TestCreateErrorOnRepositoryErrInternal(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(t, errors.Is(internalerrors.ErrInternal, err))
}

func TestGetByIdCampaignExistsCampaignSaved(t *testing.T) {
	setUp()
	repositoryMock.On("GetById", mock.MatchedBy(func(id string) bool {
		return id == campaignPending.ID
	})).Return(campaignPending, nil)

	campaignReturned, _ := service.GetById(campaignPending.ID)

	assert.Equal(t, campaignPending.ID, campaignReturned.ID)
	assert.Equal(t, campaignPending.Name, campaignReturned.Name)
	assert.Equal(t, campaignPending.Content, campaignReturned.Content)
	assert.Equal(t, campaignPending.Status, campaignReturned.Status)
	assert.Equal(t, campaignPending.CreatedBy, campaignReturned.CreatedBy)
}

func TestGetByIdErrorOnRepositoryErrInternal(t *testing.T) {
	setUp()
	repositoryMock.On("GetById", mock.Anything).Return(nil, errors.New("Something wrong'"))

	_, err := service.GetById("invalid_campaign")

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func TestDeleteCampaignWasNotFoundErrRecordNotFound(t *testing.T) {
	setUp()
	repositoryMock.On("GetById", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid_campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func TestDeleteCampaignIsNotPendingErr(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignStarted)

	err := service.Delete(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func TestDeleteErrorOnRepositoryErrInternal(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPending.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func TestDeleteCampaignWasDeletedNil(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPending == campaign
	})).Return(nil)

	err := service.Delete(campaignPending.ID)

	assert.Nil(t, err)
}

func TestStartCamapaignWasNotFoundErrRecordNotFound(t *testing.T) {
	setUp()
	repositoryMock.On("GetById", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid_campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func TestStartCampaignIsNotPendingErr(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignStarted)
	service.Repository = repositoryMock

	err := service.Start(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func TestStartCampaignWasFoundSendEmail(t *testing.T) {
	setUp()
	setUpUpdateRepository()
	setUpGetByIdRepositoryBy(campaignPending)
	emailWasSent := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignPending.ID {
			emailWasSent = true
		}
		return nil
	}
	service.SendEmail = sendMail

	service.Start(campaignPending.ID)

	assert.True(t, emailWasSent)
}

func TestStartSendEmailFailedErrInternal(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignPending)
	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendEmail = sendMail

	err := service.Start(campaignPending.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func TestStartCampaignWasUpdatedStatusIsDone(t *testing.T) {
	setUp()
	setUpSendEmailWithSuccess()
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	service.Start(campaignPending.ID)

	assert.Equal(t, campaign.Done, campaignPending.Status)
}
