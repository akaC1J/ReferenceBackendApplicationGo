package app

import (
	"log"
	"net/http"
	"route256/cart/internal/http/middleware"
	"time"

	"route256/cart/internal/app/server"
	"route256/cart/internal/http/client"
	"route256/cart/internal/pkg/repository"
	"route256/cart/internal/pkg/service/cartservice"
	"route256/cart/internal/pkg/service/productservice"
)

type App struct {
	cartRepository *repository.Repository
	httpClient     *http.Client
	productService *productservice.ProductService
	cartService    *cartservice.CartService
	router         http.Handler
}

func New(config *Config) (*App, error) {
	log.Println("[app] Starting application initialization")

	// Инициализация репозитория корзины
	cartRepository := repository.NewRepository(repository.NewStorage())

	// Создание HTTP клиента с ретраями
	httpClient := client.NewHttpClientWithRetryWithDefaultTransport(
		config.MaxRetries,
		time.Duration(config.RetryDelayMs)*time.Millisecond,
	)

	// Инициализация сервиса продуктов
	productService := productservice.NewProductService(
		httpClient,
		config.Token,
		config.ProductServiceURL,
		config.ProductServicePath,
	)

	// Инициализация сервиса корзины
	cartService := cartservice.NewService(cartRepository, productService)

	app := &App{
		cartRepository: cartRepository,
		httpClient:     httpClient,
		productService: productService,
		cartService:    cartService,
	}

	app.setupRoutes()

	return app, nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.ServeHTTP(w, r)
}

func (app *App) setupRoutes() {
	controller := server.New(app.cartService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", controller.PostItemHandleFunc)
	mux.HandleFunc("DELETE /user/{user_id}/cart", controller.DeleteCartByUserIdHandleFunc)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", controller.DeleteItemBySkuHandleFunc)
	mux.HandleFunc("GET /user/{user_id}/cart", controller.GetCartContentHandleFunc)
	app.router = middleware.NewLogMux(mux)
}
