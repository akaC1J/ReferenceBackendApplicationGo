package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) CheckoutHandleFunc(w http.ResponseWriter, r *http.Request) {
	var methodUrl = "POST /cart/checkout"
	var postCheckoutRq PostCheckoutRq

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Empty or invalid request body", methodUrl)
		return
	}

	err = json.Unmarshal(body, &postCheckoutRq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), methodUrl)
		return
	}

	orderId, err := s.cartInterface.Checkout(r.Context(), postCheckoutRq.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), methodUrl)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(PostCheckoutRs{
		orderId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to encode response", "GET /user/<user_id>/cart")
		return
	}

	w.WriteHeader(http.StatusOK)
}
