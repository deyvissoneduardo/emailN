package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/endpoints"
	"emailn/internal/infrastructure/database"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	db := database.NewBD()

	service := campaign.ServiceImpl{Repository: &database.CampaignRepository{Db: db}}
	handler := endpoints.Handler{CampaignService: &service}

	router.Post("/campaigns", endpoints.HandlerError(handler.CampaignPost))
	router.Get("/campaigns/{id}", endpoints.HandlerError(handler.CampaignGetById))

	http.ListenAndServe(":8181", router)
}
