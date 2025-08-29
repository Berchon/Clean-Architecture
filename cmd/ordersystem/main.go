package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Berchon/Clean-Architecture/configs"
	"github.com/Berchon/Clean-Architecture/internal/event/handler"
	"github.com/Berchon/Clean-Architecture/internal/infra/graph"
	"github.com/Berchon/Clean-Architecture/internal/infra/grpc/pb"
	"github.com/Berchon/Clean-Architecture/internal/infra/grpc/service"
	"github.com/Berchon/Clean-Architecture/internal/infra/web/webserver"
	"github.com/Berchon/Clean-Architecture/pkg/events"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = createOrdersTable(db)
	if err != nil {
		panic(err)
	}

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	OrderUseCase := NewOrderUseCase(db, eventDispatcher)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webCreateOrderHandler := NewWebCreateOrderHandler(db, eventDispatcher)
	webGetOrderHandler := NewWebGetOrderHandler(db, eventDispatcher)
	webserver.AddHandler(http.MethodPost, "/order", webCreateOrderHandler.Create)
	webserver.AddHandler(http.MethodGet, "/order", webGetOrderHandler.Get)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	OrderService := service.NewOrderService(*OrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, OrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		OrderUseCase: *OrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func createOrdersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id varchar(255) NOT NULL, 
		price float NOT NULL, 
		tax float NOT NULL, 
		final_price float NOT NULL, 
		PRIMARY KEY (id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
