package initialization

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	"route256/loms/internal/infra/database"
	"route256/loms/internal/infra/producer"
	grpcMW "route256/loms/internal/mw/grpc"
	httpMW "route256/loms/internal/mw/http"
	"route256/loms/internal/repository/orderrepository"
	"route256/loms/internal/repository/outboxrepository"
	"route256/loms/internal/repository/stockrepository"
	"route256/loms/internal/service/orderservice"
	"route256/loms/internal/service/processor/outboxprocessor"
	"route256/loms/internal/service/stockservice"
	transactionmanager "route256/loms/internal/service/transactionamanger"
	"route256/loms/internal/tracing"
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
	OutboxProcessor *outboxprocessor.OutboxProcessor
}

func (application *App) Run(config *Config) {
	// Создаем глобальный контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Останавливаем все процессы при завершении main функции

	// Канал для перехвата системных сигналов (Ctrl+C или завершение процесса)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Канал для ошибок сервера
	serverErrors := make(chan error, 2)

	// Запускаем gRPC сервер
	go func() {
		lis, err := net.Listen("tcp", config.GgrpcHostPort)
		if err != nil {
			serverErrors <- fmt.Errorf("failed to listen: %w", err)
			return
		}
		log.Printf("[main] gRPC server listening at %v", lis.Addr())
		if err := application.GrpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			serverErrors <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Запускаем HTTP сервер
	go func() {
		log.Printf("[main] HTTP server listening at %v", application.GwServer.Addr)
		if err := application.GwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Запускаем OutboxProcessor
	go application.OutboxProcessor.Start(ctx)

	// Ожидаем завершения или ошибки
	select {
	case sig := <-quit:
		log.Printf("[main] Caught signal %v. Shutting down gracefully...", sig)
	case err := <-serverErrors:
		log.Printf("[main] Received server error: %v", err)
	}

	// Отменяем все фоновые задачи
	cancel()

	// Завершаем работу gRPC сервера
	go func() {
		application.GrpcServer.GracefulStop()
		log.Println("[main] gRPC server stopped")
	}()

	// Завершаем работу HTTP сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := application.GwServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("[main] HTTP server forced to shutdown: %v", err)
	} else {
		log.Println("[main] HTTP server stopped")
	}

	log.Println("[main] Application gracefully stopped")
}

func MustNew(config *Config) (*App, error) {
	log.Println("[cart] Starting application initialization")

	dbRouter := database.NewDBRouter(initAllInstancesBd(config.DBConfigs))
	outboxRepository := outboxrepository.NewRepository()
	orderRepository := orderrepository.NewRepository(dbRouter, outboxRepository)
	stockRepository := stockrepository.NewRepository(dbRouter)

	stockService := stockservice.NewService(stockRepository)
	orderService := orderservice.NewService(orderRepository, stockService)

	tm := transactionmanager.NewTransactionManager(dbRouter)

	syncProducer, err := initKafkaSyncProducer(config.KafkaConfig)
	if err != nil {
		log.Fatalf("Unable to create kafka producer: %v", err)
	}
	processor := outboxprocessor.NewOutboxProcessor(outboxRepository, syncProducer, config.IntervalOutbox, tm, config.KafkaConfig.Topic)
	app := &App{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
		stockService:    stockService,
		orderService:    orderService,
		OutboxProcessor: processor,
	}

	httpMW.SwaggerUrlForCors = config.SwagerUrl
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcMW.PanicUnaryMiddleware,
			grpcMW.TracingUnaryMiddleware,
			grpcMW.MetricUnaryMiddleware,
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
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.Handle("/", gwmux)
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: httpMW.WithCorsCheckHttpHandler(mux),
	}

	app.GwServer = gwServer

	_, err = tracing.InitTracerProvider("LOMS")
	if err != nil {
		log.Fatalln(err)
	}

	return app, nil
}

func initAllInstancesBd(dbConfigs []*DBConfigs) []*database.MasterAndReplica {
	var mastersAndReplicas []*database.MasterAndReplica
	for _, dbConfig := range dbConfigs {
		dbpoolMaster, err := initDbPool(dbConfig.Master)
		if err != nil {
			log.Fatalf("Unable to create connection pool for master: %v", err)
		}

		// либо мы успешно инициализировали пул метрик, либо мы считаем что она отсутствует
		dbpoolReplica, err := initDbPool(dbConfig.ReplicaOptional)
		if err != nil {
			log.Printf("Unable to create connection pool for replice. Will be use only master mode: %v", err)
		}
		mastersAndReplicas = append(mastersAndReplicas, &database.MasterAndReplica{dbpoolMaster, dbpoolReplica})
	}
	return mastersAndReplicas

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

func initKafkaSyncProducer(kafkaConfig *KafkaConfig) (sarama.SyncProducer, error) {
	return producer.NewSyncProducer(kafkaConfig.Brokers,
		producer.WithMaxRetries(int(kafkaConfig.RetryMax)),
		producer.WithRetryBackoff(500*time.Millisecond))
}
