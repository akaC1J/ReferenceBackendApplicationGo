package main

import (
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
)

func main() {

	log.Println("[main] Application start initialization")
	loadEnv()
	cartRepository := repository.NewRepository()

	maxRetries, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if err != nil {
		panic(err)
	}

	retryMs, err := strconv.Atoi(os.Getenv("RETRY_DELAY_MS"))
	if err != nil {
		panic(err)
	}

	httpClient := client.NewHttpClientWithRetry(maxRetries, int64(retryMs))

	//Внимание это чувствительная информация, на проде надо реализовать получения секрета по-другому
	// например через hashicorp vault
	//но в качестве теста можно
	token := os.Getenv("TOKEN")
	baseUrl := os.Getenv("PRODUCT_SERVICE_BASE_URL")
	path := os.Getenv("PRODUCT_SERVICE_PATH")
	productService := productservice.NewProductService(httpClient, token, baseUrl, path)

	cartService := cartservice.NewService(cartRepository, productService)

	controller := server.New(cartService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", controller.PostItemHandleFunc)
	mux.HandleFunc("DELETE /user/{user_id}/cart", controller.DeleteCartByUserIdHandleFunc)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", controller.DeleteItemBySkuHandleFunc)
	mux.HandleFunc("GET /user/{user_id}/cart", controller.GetCartContentHandleFunc)

	logMux := middleware.NewLogMux(mux)

	log.Println("[main] Application initialization successes")

	serverHost := os.Getenv("HOST")
	serverPort := os.Getenv("PORT")
	connectParam := serverHost + ":" + serverPort
	log.Printf("[main] Starting server on %s\n", connectParam)
	if err := http.ListenAndServe(connectParam, logMux); err != nil {
		panic(err)
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
