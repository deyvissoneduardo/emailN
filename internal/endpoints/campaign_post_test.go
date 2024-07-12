package endpoints

import (
	"bytes"
	"context"
	"emailn/internal/contracts"
	internalmock "emailn/internal/test/internal_mock"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup(body contracts.NewCampaignDto, createdByExpected string) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("POST", "/", &buf)
	ctx := context.WithValue(req.Context(), "email", createdByExpected)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	return req, rr
}

func TestCampaignsPostShouldSaveNewCamapaign(t *testing.T) {
	assert := assert.New(t)
	createdByExpected := "teste1@teste.com.br"
	body := contracts.NewCampaignDto{
		Name:    "teste",
		Content: "Hi everyone",
		Emails:  []string{"teste@teste.com"},
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.MatchedBy(func(request contracts.NewCampaignDto) bool {
		if request.Name == body.Name &&
			request.Content == body.Content &&
			request.CreateBy == createdByExpected {
			return true
		} else {
			return false
		}
	})).Return("34x", nil)
	handler := Handler{CampaignService: service}
	req, rr := setup(body, createdByExpected)

	_, status, err := handler.CampaignPost(rr, req)

	assert.Equal(201, status)
	assert.Nil(err)
}

func TestCampaignsPostShouldInformErrorWhenExist(t *testing.T) {
	assert := assert.New(t)
	body := contracts.NewCampaignDto{
		Name:    "teste",
		Content: "Hi everyone",
		Emails:  []string{"teste@teste.com"},
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	handler := Handler{CampaignService: service}
	req, rr := setup(body, "teste@teste.com.br")

	_, _, err := handler.CampaignPost(rr, req)

	assert.NotNil(err)
}
