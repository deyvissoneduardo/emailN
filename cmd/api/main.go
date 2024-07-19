package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/endpoints"
	"emailn/internal/infrastructure/database"
	"emailn/internal/infrastructure/mail"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	errMail := mail.SendMail()
	if errMail != nil {
		log.Fatal(errMail.Error())
	}

	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	db := database.NewBD()

	service := campaign.ServiceImpl{Repository: &database.CampaignRepository{Db: db}}
	handler := endpoints.Handler{CampaignService: &service}

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.Route("/campaigns", func(router chi.Router) {
		router.Use(endpoints.Auth)
		router.Post("/", endpoints.HandlerError(handler.CampaignPost))
		router.Get("/{id}", endpoints.HandlerError(handler.CampaignGetById))
		router.Delete("/delete/{id}", endpoints.HandlerError(handler.CampaignDelete))
	})

	http.ListenAndServe(":8181", router)
}
