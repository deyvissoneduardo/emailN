package mock

import (
	"emailn/internal/contracts"

	"github.com/stretchr/testify/mock"
)

type CampaignServiceMock struct {
	mock.Mock
}

func (r *CampaignServiceMock) Create(newCampaign contracts.NewCampaignDto) (string, error) {
	args := r.Called(newCampaign)
	return args.String(0), args.Error(1)
}

func (r *CampaignServiceMock) GetById(id string) (*contracts.CampaignResponse, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contracts.CampaignResponse), args.Error(1)
}

func (r *CampaignServiceMock) Delete(id string) error {
	return nil
}
