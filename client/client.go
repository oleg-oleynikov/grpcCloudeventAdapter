package main

import (
	"context"
	"grpccloudevents"
)

func main() {
	client := grpccloudevents.NewCloudeventsServiceClient("localhost:50051")
	client.Subscribe(context.Background(), "example.com")
}
