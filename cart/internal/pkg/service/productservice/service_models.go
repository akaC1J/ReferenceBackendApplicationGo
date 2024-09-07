package productservice

type productGetProductRequest struct {
	Token string `json:"token"`
	Sku   string `json:"sku"`
}

type productGetProductResponse struct {
	Name  string `json:"name" validate:"required"`
	Price uint32 `json:"price" validate:"required,gt=0"`
}
