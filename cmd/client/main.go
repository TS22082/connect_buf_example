package main

import (
	"context"
	"fmt"
	"net/http"

	greetv1 "github.com/TS22082/connect_buf_example/gen/greet/v1"
	"github.com/TS22082/connect_buf_example/gen/greet/v1/greetv1connect"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	res, err := client.Greet(
		context.Background(),
		&greetv1.GreetRequest{Name: "Billy Joel"},
	)

	if err != nil {
		fmt.Println("FUUUCK")
	}

	fmt.Println(res.Greeting)
}
