package endpoints

import (
	"emailn/internal/contracts"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCampaignsGetByIdShouldReturnCampaign(t *testing.T) {
	setup()
	campaignId := "343"
	campaign := contracts.CampaignResponse{
		ID:      campaignId,
		Name:    "Test",
		Content: "Hi!",
		Status:  "Pending",
	}
	service.On("GetById", mock.Anything).Return(&campaign, nil)
	req, rr := newHttpTest("GET", "/", nil)
	req = addParameter(req, "id", campaignId)

	response, status, _ := handler.CampaignGetById(rr, req)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, campaign.ID, response.(*contracts.CampaignResponse).ID)
	assert.Equal(t, campaign.Name, response.(*contracts.CampaignResponse).Name)
}

func TestCampaignsGetByIdShouldReturnErrorWhenSomethingWrong(t *testing.T) {
	setup()
	errExpected := errors.New("something wrong")
	service.On("GetById", mock.Anything).Return(nil, errExpected)
	req, rr := newHttpTest("GET", "/", nil)

	_, _, errReturned := handler.CampaignGetById(rr, req)

	assert.Equal(t, errExpected.Error(), errReturned.Error())
}
