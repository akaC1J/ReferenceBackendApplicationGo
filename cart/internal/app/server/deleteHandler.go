package server

import (
	"net/http"
	"route256/cart/internal/pkg/model"
	"strconv"
)

func (s *Server) DeleteItemBySkuHandleFunc(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	rawSkuId := r.PathValue("sku_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id format", "DELETE /user/<user_id>/cart/<sku_id>")
		return
	}

	skuId, err := strconv.ParseInt(rawSkuId, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid sku_id format", "DELETE /user/<user_id>/cart/<sku_id>")
		return
	}

	err = s.reviewService.DeleteCartItem(r.Context(), model.UserId(userId), model.SKU(skuId))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "DELETE /user/<user_id>/cart/<sku_id>")
	}

	w.WriteHeader(204)
}

func (s *Server) DeleteCartByUserIdHandleFunc(w http.ResponseWriter, r *http.Request) {
	rawUserId := r.PathValue("user_id")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id format", "DELETE /user/<user_id>/cart")
		return
	}

	err = s.reviewService.CleanUpCart(r.Context(), model.UserId(userId))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "DELETE /user/<user_id>/cart")
		return
	}

	w.WriteHeader(204)
}
