package contracts

type CampaignResponse struct {
	ID                   string
	Name                 string
	Content              string
	Status               string
	CreatedBy            string
	AmountOfEmailsToSend int
}
