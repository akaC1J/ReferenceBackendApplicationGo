package initialization

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/loms/internal/app/grpccontroller"
	"route256/loms/internal/generated/api/loms/v1"
	"route256/loms/internal/infra"
	grpcMW "route256/loms/internal/mw/grpc"
	httpMW "route256/loms/internal/mw/http"
	"route256/loms/internal/repository/orderrepository"
	"route256/loms/internal/repository/stockrepository"
	"route256/loms/internal/service/orderservice"
	"route256/loms/internal/service/stockservice"
	"syscall"
	"time"
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
	log.Printf("[main] gRPC server listening at %v", lis.Addr())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 2)

	go func() {
		if err := application.GrpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Printf("[main] gRPC server error: %v", err)
			serverErrors <- err
		}
	}()

	log.Printf("[main] HTTP server listening at %v", application.GwServer.Addr)

	go func() {
		if err := application.GwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[main] HTTP server error: %v", err)
			serverErrors <- err
		}
	}()

	select {
	case sig := <-quit:
		log.Printf("[main] Caught signal %v. Shutting down gracefully...", sig)
	case err := <-serverErrors:
		log.Printf("[main] Received server error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		application.GrpcServer.GracefulStop()
		log.Println("[main] gRPC server stopped")
	}()

	if err := application.GwServer.Shutdown(ctx); err != nil {
		log.Printf("[main] HTTP server forced to shutdown: %v", err)
	} else {
		log.Println("[main] HTTP server stopped")
	}

	log.Println("[main] Application gracefully stopped")
}

func MustNew(config *Config) (*App, error) {
	log.Println("[cart] Starting application initialization")
	dbpoolMaster, err := initDbPool(config.DBConfigMaster)
	if err != nil {
		log.Fatalf("Unable to create connection pool for master: %v", err)
	}

	// либо мы успешно инициализировали пул метрик, либо мы считаем что она отсутствует
	dbpoolReplica, err := initDbPool(config.DBConfigReplicaOptional)
	if err != nil {
		log.Printf("Unable to create connection pool for replice. Will be use only master mode: %v", err)
	}
	dbRouter := infra.NewDBRouter(dbpoolMaster, dbpoolReplica)

	orderRepository := orderrepository.NewRepository(dbRouter)
	stockRepository := stockrepository.NewRepository(dbRouter)

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
	loms.RegisterLomsServer(grpcServer, lomsController)

	app.GrpcServer = grpcServer

	conn, err := grpc.NewClient(config.GgrpcHostPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to deal:", err)
	}

	gwmux := runtime.NewServeMux()

	if err = loms.RegisterLomsHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: grpcMW.WithCorsCheckHttpHandler(httpMW.WithHTTPLoggingMiddleware(gwmux)),
	}

	app.GwServer = gwServer

	return app, nil
}

func initDbPool(dbConfig *DBConfig) (*pgxpool.Pool, error) {
	if dbConfig == nil {
		return nil, fmt.Errorf("dbConfig is nil")
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s",
		dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHostPort, dbConfig.DBName)
	poolConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("unable to parse configuration for db (host:port %s): %w", dbConfig.DBHostPort, err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create pool with host:port %s: %w", dbConfig.DBHostPort, err)
	}
	return pool, nil
}
