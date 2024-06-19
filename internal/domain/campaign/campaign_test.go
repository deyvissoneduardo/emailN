package campaign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCampaign(t *testing.T) {
	assert := assert.New(t)

	var (
		name     = "Campaign X"
		content  = "Body Hi!"
		contacts = []string{"email1@e.com", "email2@e.com"}
		// createdBy = "teste@teste.com.br"
		// fake      = faker.New()
	)
	campaign := NewCampaign(name, content, contacts)

	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))

}
