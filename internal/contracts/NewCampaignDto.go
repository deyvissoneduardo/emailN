package contracts

type NewCampaignDto struct {
	Name     string
	Content  string
	Emails   []string
	CreateBy string
}
