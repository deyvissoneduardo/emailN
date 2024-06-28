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

	service := campaign.ServiceImpl{Repository: &database.CampaignRepository{}}
	handler := endpoints.Handler{CampaignService: &service}

	router.Post("/campaigns", endpoints.HandlerError(handler.CampaignPost))
	router.Get("/campaigns", endpoints.HandlerError(handler.CampaignGet))

	http.ListenAndServe(":8181", router)
}
