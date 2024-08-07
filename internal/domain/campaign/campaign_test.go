package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var (
	name      = "Campaign X"
	content   = "Body Hi!"
	contacts  = []string{"email1@e.com", "email2@e.com"}
	createdBy = "test@gmail.com"
	fake      = faker.New()
)

func TestNewCampaignCreateCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
	assert.Equal(campaign.CreatedBy, createdBy)
}

func TestNewCampaignIDIsNotNil(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.NotNil(campaign.ID)
}

func TestNewCampaignCreatedMustBeNow(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Greater(campaign.CreatedOn, now)
}

func TestNewCampaignMustValidateMinName(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("", content, contacts, createdBy)

	assert.Equal("name is required with min 5", err.Error())
}

func TestNewCampaignMustValidateMaxName(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts, createdBy)

	assert.Equal("name is required with max 24", err.Error())
}

func TestNewCampaignMustValidateContent(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, "", contacts, createdBy)

	assert.Equal("content is required with min 5", err.Error())
}

func TestNewCampaignMustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, nil, createdBy)

	assert.Equal("contacts is required with min 1", err.Error())
}

func TestNewCampaignMustStatusStartWithPending(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Equal(Pending, campaign.Status)
}
