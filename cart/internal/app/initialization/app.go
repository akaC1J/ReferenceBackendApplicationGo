package initialization

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"route256/cart/internal/generated/api/loms/v1"
	"route256/cart/internal/http/middleware"
	"route256/cart/internal/pkg/service/lomsservice"
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
	log.Println("[cart] Starting application initialization")

	// Инициализация репозитория корзины
	cartRepository := repository.NewRepository(repository.NewStorage())

	// Создание HTTP клиента с ретраями
	httpClient := client.NewHttpClientWithRetryWithDefaultTransport(
		config.MaxRetries,
		time.Duration(config.RetryDelayMs)*time.Millisecond,
	)

	limitTripper := client.NewLimiterRoundTripper(httpClient.Transport, config.RequestsPerSecond)

	grpcClient := newGRPCClient(config)

	// Инициализация сервиса заказов
	lomsService := lomsservice.NewLomsService(grpcClient)

	// Инициализация сервиса продуктов
	productService := productservice.NewProductService(
		limitTripper,
		config.Token,
		config.ProductServiceURL,
		config.ProductServicePath,
	)

	// Инициализация сервиса корзины
	cartService := cartservice.NewService(cartRepository, productService, lomsService)

	app := &App{
		cartRepository: cartRepository,
		httpClient:     httpClient,
		productService: productService,
		cartService:    cartService,
	}

	app.setupRoutes()

	return app, nil
}

func newGRPCClient(config *Config) loms.LomsClient {

	conn, err := grpc.NewClient("dns:///"+config.LomsBaseUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil
	}
	return loms.NewLomsClient(conn)
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
	mux.HandleFunc("POST /cart/checkout", controller.CheckoutHandleFunc)
	app.router = middleware.NewLogMux(mux)
}
