package server

import (
	"encoding/json"
	"io"
	"net/http"
	"route256/cart/internal/pkg/model"
)

func (s *Server) PostItemHandleFunc(w http.ResponseWriter, r *http.Request) {
	userId, err := getParamFromReq(r, "user_id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "DELETE /user/<user_id>/cart/<sku_id>")
		return
	}

	skuId, err := getParamFromReq(r, "sku_id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "DELETE /user/<user_id>/cart/<sku_id>")
		return
	}

	var postItemRq PostItemRequest

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Empty or invalid request body", "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	err = json.Unmarshal(body, &postItemRq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	_, err = s.cartInterface.AddCartItem(r.Context(), model.CartItem{
		SKU:    model.SKU(skuId),
		UserId: model.UserId(userId),
		Count:  postItemRq.Count,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	w.WriteHeader(http.StatusOK)
}
