package initialization

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	"os"
	"route256/loms/internal/app/grpccontroller"
	"route256/loms/internal/model"
	"route256/loms/internal/mw"
	"route256/loms/internal/repository/orderrepository"
	"route256/loms/internal/repository/stockrepository"
	"route256/loms/internal/service/orderservice"
	"route256/loms/internal/service/stockservice"
	lomsGrpc "route256/loms/pkg/api/loms/v1"
)

type App struct {
	orderRepository *orderrepository.Repository
	stockRepository *stockrepository.Repository
	stockService    *stockservice.Service
	orderService    *orderservice.Service
	GrpcServer      *grpc.Server
	GwServer        *http.Server
}

func New(config *Config) (*App, error) {
	log.Println("[cart] Starting application initialization")

	orderRepository := orderrepository.NewRepository()
	stockRepository := initStockRepositoryFromFile(config.StockFilePath)

	stockService := stockservice.NewService(stockRepository)
	orderService := orderservice.NewService(orderRepository, stockService)

	app := &App{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
		stockService:    stockService,
		orderService:    orderService,
	}

	mw.SwaggerUrlForCors = config.SwagerUrl
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Panic,
			mw.Logger,
			mw.Validate,
		),
	)
	reflection.Register(grpcServer)
	lomsGrpc.RegisterLomsServer(grpcServer, grpccontroller.NewLomsController(app.orderService, app.stockService))

	app.GrpcServer = grpcServer

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to deal:", err)
	}

	gwmux := runtime.NewServeMux()

	if err = lomsGrpc.RegisterLomsHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: mw.WithCorsCheckHttpHandler(mw.WithHTTPLoggingMiddleware(gwmux)),
	}

	app.GwServer = gwServer

	return app, nil
}

func initStockRepositoryFromFile(pathToStockDataFile string) *stockrepository.Repository {
	open, err := os.Open(pathToStockDataFile)
	if err != nil {
		log.Fatalf("[cart] Error opening stock data file: %v", err)
	}
	defer open.Close()
	var stocks []*model.Stock
	err = json.NewDecoder(open).Decode(&stocks)
	if err != nil {
		log.Fatalf("[cart] Error decoding stock data file: %v", err)
	}
	return stockrepository.NewRepository(stocks)
}
