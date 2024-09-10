package server

import (
	"net/http"
	"route256/cart/internal/pkg/model"
)

func (s *Server) DeleteItemBySkuHandleFunc(w http.ResponseWriter, r *http.Request) {
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

	err = s.cartInterface.DeleteCartItem(r.Context(), model.UserId(userId), model.SKU(skuId))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "DELETE /user/<user_id>/cart/<sku_id>")
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DeleteCartByUserIdHandleFunc(w http.ResponseWriter, r *http.Request) {
	userId, err := getParamFromReq(r, "user_id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "DELETE /user/<user_id>/cart/<sku_id>")
		return
	}

	err = s.cartInterface.CleanUpCart(r.Context(), model.UserId(userId))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), "DELETE /user/<user_id>/cart")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
