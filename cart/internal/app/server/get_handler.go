package server

import (
	"encoding/json"
	"net/http"
	"route256/cart/internal/pkg/model"
)

func (s *Server) GetCartContentHandleFunc(w http.ResponseWriter, r *http.Request) {
	userId, err := getParamFromReq(r, "user_id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "GET /user/<user_id>/cart")
		return
	}

	cartContent, err := s.cartInterface.GetCartItem(r.Context(), model.UserId(userId))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "GET /user/<user_id>/cart")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(GetCartContentResponse{
		cartContent,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to encode response", "GET /user/<user_id>/cart")
		return
	}
}
