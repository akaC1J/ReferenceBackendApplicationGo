package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"route256/cart/internal/app/server"
	"route256/cart/internal/http/client"
	"route256/cart/internal/http/middleware"
	"route256/cart/internal/pkg/repository"
	"route256/cart/internal/pkg/service/cartservice"
	"route256/cart/internal/pkg/service/productservice"
	"strconv"
	"strings"
	"time"
)

type App struct {
	cartRepository *repository.Repository
	httpClient     *http.Client
	productService *productservice.ProductService
	cartService    *cartservice.CartService
}

func NewApp() (*App, error) {
	log.Println("[main] Application start initialization")
	loadEnv()

	cartRepository := repository.NewRepository(repository.NewStorage())

	maxRetries, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse MAX_RETRIES: %w", err)
	}

	retryMs, err := strconv.Atoi(os.Getenv("RETRY_DELAY_MS"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RETRY_DELAY_MS: %w", err)
	}

	httpClient := client.NewHttpClientWithRetryWithDefaultTransport(maxRetries, time.Duration(retryMs)*time.Millisecond)

	token := os.Getenv("TOKEN")
	baseUrl := os.Getenv("PRODUCT_SERVICE_BASE_URL")
	path := os.Getenv("PRODUCT_SERVICE_PATH")
	productService := productservice.NewProductService(httpClient, token, baseUrl, path)

	cartService := cartservice.NewService(cartRepository, productService)

	return &App{
		cartRepository: cartRepository,
		httpClient:     httpClient,
		productService: productService,
		cartService:    cartService,
	}, nil
}

func (app *App) SetupRoutes() *http.ServeMux {
	controller := server.New(app.cartService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", controller.PostItemHandleFunc)
	mux.HandleFunc("DELETE /user/{user_id}/cart", controller.DeleteCartByUserIdHandleFunc)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", controller.DeleteItemBySkuHandleFunc)
	mux.HandleFunc("GET /user/{user_id}/cart", controller.GetCartContentHandleFunc)
	return mux
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("[main] Failed to initialize application: %v", err)
	}

	mux := app.SetupRoutes()
	logMux := middleware.NewLogMux(mux)

	log.Println("[main] Application initialization successful")

	serverAddress := os.Getenv("HOST") + ":" + os.Getenv("PORT")
	log.Printf("[main] Starting server on %s\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, logMux); err != nil {
		log.Fatalf("[main] Failed to start server: %v", err)
	}
}

func loadEnv() {
	// Загружаем файл .env, если ENV=test, загружаем .env.test
	env := os.Getenv("ENV")
	if strings.ToUpper(env) == "TEST" {
		log.Println("[main] Received ENV=test. Load .env.test in environment")
		err := godotenv.Load(".env.test")

		if err != nil {
			panic(err)
		}
	} else {
		log.Println("[main] Load .env in environment")
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}
}
