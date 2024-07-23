package endpoints

import (
	"emailn/internal/contracts"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	createdByExpected = "teste1@teste.com.br"
	body              = contracts.NewCampaignDto{
		Name:    "teste",
		Content: "Hi everyone",
		Emails:  []string{"teste@teste.com"},
	}
)

func TestCampaignsPostShouldSaveNewCamapaign(t *testing.T) {
	setup()
	service.On("Create", mock.MatchedBy(func(request contracts.NewCampaignDto) bool {
		if request.Name == body.Name &&
			request.Content == body.Content &&
			request.CreateBy == createdByExpected {
			return true
		} else {
			return false
		}
	})).Return("34x", nil)
	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, status, err := handler.CampaignPost(rr, req)

	assert.Equal(t, http.StatusOK, status)
	assert.Nil(t, err)
}

func TestCampaignsPostShouldInformErrorWhenExist(t *testing.T) {
	setup()
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, _, err := handler.CampaignPost(rr, req)

	assert.NotNil(t, err)
}
