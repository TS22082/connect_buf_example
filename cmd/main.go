package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	greetv1 "github.com/TS22082/connect_buf_example/gen/greet/v1"
	"github.com/TS22082/connect_buf_example/gen/greet/v1/greetv1connect"
	"github.com/TS22082/connect_buf_example/internal"
	"github.com/gorilla/mux"
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	_ context.Context,
	req *greetv1.GreetRequest,
) (*greetv1.GreetResponse, error) {
	res := &greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Name),
	}
	return res, nil
}

func main() {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/user/{id}", internal.TestHandler).Methods("GET")

	greeter := &GreetServer{}
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
