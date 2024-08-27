package main

import (
	"grpccloudevents"
	"grpccloudevents/proto/resume"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	v1 "open-cluster-management.io/sdk-go/pkg/cloudevents/generic/options/grpc/protobuf/v1"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	resumeEvent := &resume.ResumeCreated{
		ResumeId: "12121",
		UserName: "user1",
		Password: "pypa",
	}

	protoEvent, err := anypb.New(resumeEvent)
	if err != nil {
		log.Println("Failed to cast to anypb")
		return
	}

	event := &v1.CloudEvent{
		Id:          "1231312345",
		Source:      "example.com",
		SpecVersion: "1.0",
		Type:        "resume.send",
		Data: &v1.CloudEvent_ProtoData{
			ProtoData: protoEvent,
		},
	}

	client := grpccloudevents.NewCloudeventsServiceClient("localhost:50051")

	startTime := time.Now()

	for i := 0; i < 10000; i++ {
		client.Publish(event)
	}

	log.Println("Затрачено: ", time.Since(startTime))
}
