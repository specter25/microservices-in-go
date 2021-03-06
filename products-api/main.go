package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	protos "github.com/specter25/microservices-in-go/currency-api/protos/currency"
	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	"github.com/specter25/microservices-in-go/products-api/data"
	"github.com/specter25/microservices-in-go/products-api/handlers"
)

func main() {

	//create a currency client to the currency service

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cc := protos.NewCurrencyClient(conn)

	l := hclog.Default()
	v := data.NewValidation()

	//create databse insatnce
	db := data.NewProductDB(cc, l)

	ph := handlers.NewProducts(l, v, db)

	// gh := handlers.NewGoodbye(l)
	//create a new swerve mux
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts).Queries("currency", "{[A_Z]{3}}")
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle).Queries("currency", "{[A_Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)

	//Redoc middleware used to render the swagger documentation properly study redoc documentation to understand this code
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)

	//we have to create a file server to serve the file that we want from a server
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// go func used so that it does not stop running
	go func() {
		l.Info("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	//now everything gets blocked here untill the message is received
	sig := <-sigChan
	l.Error("Received terminate , graceful shutdown", sig)
	//This is very important to shutdown the server after it has received all the request

	//shutdown needs a context
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
