package initialization

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"net/http/pprof"
	"route256/cart/internal/generated/api/loms/v1"
	"route256/cart/internal/http/middleware"
	"route256/cart/internal/logger"
	"route256/cart/internal/metrics"
	"route256/cart/internal/pkg/service/lomsservice"
	"route256/cart/internal/tracing"
	"time"

	_ "net/http/pprof"
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
	Logger         *logger.Logger
}

func New(config *Config, pkgLoger *logger.Logger) (*App, error) {
	log.Println("[cart] Starting application initialization")

	// Инициализация репозитория корзины
	cartRepository := repository.NewRepository(repository.NewStorage())

	// Создание HTTP клиента с ретраями
	httpClient := client.NewHttpClientWithRetryWithDefaultTransport(
		config.MaxRetries,
		time.Duration(config.RetryDelayMs)*time.Millisecond,
	)

	limitTripper := client.NewLimiterRoundTripper(httpClient.Transport, config.RequestsPerSecond)
	metricTripper := client.NewMetricTripper(limitTripper)
	grpcClient := newGRPCClient(config)

	// Инициализация сервиса заказов
	lomsService := lomsservice.NewLomsService(grpcClient)

	// Инициализация сервиса продуктов
	productService := productservice.NewProductService(
		metricTripper,
		config.Token,
		config.ProductServiceURL,
		config.ProductServicePath,
	)

	cacheProductService := productservice.NewCacheProductService(config.CacheCapacity, productService)

	// Инициализация сервиса корзины
	cartService := cartservice.NewService(cartRepository, cacheProductService, lomsService)

	app := &App{
		cartRepository: cartRepository,
		httpClient:     httpClient,
		productService: productService,
		cartService:    cartService,
		Logger:         pkgLoger,
	}

	_, err := tracing.InitTracerProvider("CART")
	if err != nil {
		return nil, err
	}

	app.setupRoutes()

	return app, nil
}

type metadataCarrier metadata.MD

func (mc metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range mc {
		keys = append(keys, k)
	}
	return keys
}

func (mc metadataCarrier) Get(key string) string {
	values := metadata.MD(mc).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (mc metadataCarrier) Set(key string, value string) {
	metadata.MD(mc).Append(key, value)
}

func newGRPCClient(config *Config) loms.LomsClient {
	unaryInterceptor := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Извлекаем текущий контекст трейсинга и добавляем его в метаданные
		md := metadata.New(nil)
		otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(md))

		// Обновляем контекст с метаданными
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Запуск gRPC-запроса и измерение времени выполнения
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		// Запись метрик
		metrics.RecordExternalRequest(method, status.Code(err).String(), duration)

		return err
	}

	// Подключение gRPC клиента с интерсептором
	conn, err := grpc.Dial("dns:///"+config.LomsBaseUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryInterceptor))
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

	buisnessMux := middleware.NewLogMux(middleware.NewTraceMux(middleware.NewMetricMux(mux)))

	mux = http.NewServeMux()
	mux.Handle("/", buisnessMux)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	app.router = mux

}
