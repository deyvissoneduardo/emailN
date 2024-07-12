package endpoints

import (
	"emailn/internal/contracts"
	"net/http"

	"github.com/go-chi/render"
)

func (h *Handler) CampaignPost(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var request contracts.NewCampaignDto

	render.DecodeJSON(r.Body, &request)

	email := r.Context().Value("email").(string)
	request.CreateBy = email

	id, err := h.CampaignService.Create(request)

	return map[string]string{"id": id}, http.StatusCreated, err
}
