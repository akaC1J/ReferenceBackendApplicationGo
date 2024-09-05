package server

import (
	"encoding/json"
	"io"
	"net/http"
	"route256/cart/internal/pkg/model"
	"strconv"
)

type PostItemRequest struct {
	Count uint16 `json:"count"`
}

func (s *Server) PostItemHandleFunc(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	rawSkuId := r.PathValue("sku_id")

	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id format", "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	skuId, err := strconv.ParseInt(rawSkuId, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid sku_id format", "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	var postItemRq PostItemRequest

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	err = json.Unmarshal(body, &postItemRq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	_, err = s.reviewService.AddCartItem(r.Context(), model.CartItem{
		SKU:    model.SKU(skuId),
		UserId: model.UserId(userId),
		Count:  postItemRq.Count,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "POST /user/<user_id>/cart/<sku_id>")
		return
	}

	w.WriteHeader(200)
}
