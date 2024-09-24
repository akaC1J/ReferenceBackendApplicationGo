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
	"net"
	"net/http"
	"os"
	"route256/loms/internal/app/grpccontroller"
	"route256/loms/internal/model"
	grpcMW "route256/loms/internal/mw/grpc"
	httpMW "route256/loms/internal/mw/http"
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

func (application *App) Run(config *Config) {
	lis, err := net.Listen("tcp", config.GgrpcHostPort)
	if err != nil {
		log.Fatalf("[main] failed to listen: %v", err)
	}

	log.Println("[main] Application initialization successful")
	log.Printf("[main] server listening at %v", lis.Addr())

	go func() {
		if err = application.GrpcServer.Serve(lis); err != nil {
			log.Fatalf("[main] failed to serve: %v", err)
		}
	}()
	log.Printf("[main] server listening at %v", application.GwServer.Addr)
	if err = application.GwServer.ListenAndServe(); err != nil {
		log.Fatalf("[main] failed to serve: %v", err)
	}
}

func MustNew(config *Config) (*App, error) {
	log.Println("[cart] Starting application initialization")

	orderRepository := orderrepository.NewRepository()
	stockRepository := mustNewStockRepositoryFromFile(config.StockFilePath)

	stockService := stockservice.NewService(stockRepository)
	orderService := orderservice.NewService(orderRepository, stockService)

	app := &App{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
		stockService:    stockService,
		orderService:    orderService,
	}

	grpcMW.SwaggerUrlForCors = config.SwagerUrl
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcMW.PanicUnaryMiddleware,
			grpcMW.LoggingUnaryMiddleware,
			grpcMW.ValidateUnaryMiddleware,
		),
	)
	reflection.Register(grpcServer)
	lomsController := grpccontroller.NewLomsController(app.orderService, app.stockService)
	lomsGrpc.RegisterLomsServer(grpcServer, lomsController)

	app.GrpcServer = grpcServer

	conn, err := grpc.NewClient(config.GgrpcHostPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to deal:", err)
	}

	gwmux := runtime.NewServeMux()

	if err = lomsGrpc.RegisterLomsHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: grpcMW.WithCorsCheckHttpHandler(httpMW.WithHTTPLoggingMiddleware(gwmux)),
	}

	app.GwServer = gwServer

	return app, nil
}

func mustNewStockRepositoryFromFile(pathToStockDataFile string) *stockrepository.Repository {
	open, err := os.Open(pathToStockDataFile)
	if err != nil {
		log.Fatalf("[cart] Error opening stock data file: %v", err)
	}
	defer open.Close()
	var stocksFromFile []*model.Stock
	err = json.NewDecoder(open).Decode(&stocksFromFile)
	var mapStock = make(map[model.SKUType]model.Stock)
	for _, stock := range stocksFromFile {
		mapStock[stock.SKU] = *stock
	}
	if err != nil {
		log.Fatalf("[cart] Error decoding stock data file: %v", err)
	}
	return stockrepository.NewRepository(mapStock)
}
