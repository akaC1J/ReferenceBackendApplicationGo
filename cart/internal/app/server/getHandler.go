package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cartservice"
	"strconv"
)

type GetCartContentResponse struct {
	*cartservice.CartContent
}

func (s *Server) GetCartContentHandleFunc(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id format", "GET /user/<user_id>/cart")
		return
	}

	cartContent, err := s.reviewService.GetCartItem(r.Context(), model.UserId(userId))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "GET /user/<user_id>/cart")
		return
	}

	rawResponse, err := json.Marshal(GetCartContentResponse{
		cartContent,
	})

	_, err = fmt.Fprint(w, string(rawResponse))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), "GET /user/<user_id>/cart")
		return
	}
	w.WriteHeader(200)
}
