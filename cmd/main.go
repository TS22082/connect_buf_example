package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/TS22082/connect_buf_example/gen/greet/v1/greetv1connect"
	"github.com/TS22082/connect_buf_example/internal/handlers"
	"github.com/TS22082/connect_buf_example/internal/services/greet"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/user/{id}", handlers.TestHandler).Methods("GET")

	greeter := &greet.GreetServer{}
	path, handler := greetv1connect.NewGreetServiceHandler(
		greeter,
		connect.WithInterceptors(validate.NewInterceptor()),
	)

	router.PathPrefix(path).Handler(handler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	s := http.Server{
		Addr:      "localhost:8080",
		Handler:   router,
		Protocols: p,
	}

	fmt.Println("REST route: GET /api/v1/user/{id}")
	fmt.Printf("RPC route: POST %s{method}\n", path)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed to start: %v", err)
	}

}
