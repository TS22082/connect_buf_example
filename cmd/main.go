package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	greetv1 "github.com/TS22082/connect_buf_example/gen/greet/v1"
	"github.com/TS22082/connect_buf_example/gen/greet/v1/greetv1connect"
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
	api.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectId, exists := vars["id"]

		if !exists {
			http.Error(w, "Wrong input", http.StatusExpectationFailed)
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"msg":       fmt.Sprintf("ProjectId: %s", projectId),
			"timestamp": time.Now().UTC(),
		}); err != nil {
			http.Error(w, "Can not send response", http.StatusExpectationFailed)
			return
		}
	}).Methods("GET")

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
