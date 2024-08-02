package grpccloudevents

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "open-cluster-management.io/sdk-go/pkg/cloudevents/generic/options/grpc/protobuf/v1"
)

var (
	counterEvents = 0
)

type CloudeventsServiceClient struct {
	sync.RWMutex
	// eventChan  chan *v1.CloudEvent
	// errChan    chan error
	grpcClient v1.CloudEventServiceClient
	wg         sync.WaitGroup
}

func NewCloudeventsServiceClient(target string) *CloudeventsServiceClient {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil
	}
	grpcClient := v1.NewCloudEventServiceClient(conn)

	c := &CloudeventsServiceClient{
		// eventChan:  make(chan *v1.CloudEvent, 100),
		// errChan:    make(chan error, 10),
		grpcClient: grpcClient,
	}

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	c.processPublish()
	// }()

	return c
}

func (c *CloudeventsServiceClient) Publish(cloudEvent *v1.CloudEvent) error {
	_, err := c.grpcClient.Publish(context.TODO(), &v1.PublishRequest{
		Event: cloudEvent,
	})
	if err != nil {
		return err
	}
	return nil
	// СИНХРОННО

	// c.wg.Add(1)
	// go func() {
	// 	defer c.wg.Done()
	// 	_, err := c.grpcClient.Publish(context.TODO(), &v1.PublishRequest{
	// 		Event: cloudEvent,
	// 	})
	// 	if err != nil {
	// 		log.Printf("Error publishing event: %v", err)
	// 	}
	// }()
	// return nil
}

// func (c *CloudeventsServiceClient) processPublish() {
// 	for event := range c.eventChan {
// 		backoff := 100 * time.Millisecond
// 		maxBackoff := 5 * time.Second
// 		for attempt := 0; attempt < 5; attempt++ {
// 			_, err := c.grpcClient.Publish(context.TODO(), &v1.PublishRequest{
// 				Event: event,
// 			})
// 			if err == nil {
// 				break
// 			}
// 			log.Printf("Error when publishing event: %s, retrying in %s", err, backoff)
// 			time.Sleep(backoff)
// 			backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
// 		}
// 	}
// }

func (c *CloudeventsServiceClient) Subscribe(ctx context.Context, source string) {
	subReq := &v1.SubscriptionRequest{
		Source: source,
	}

	stream, err := c.grpcClient.Subscribe(ctx, subReq)
	if err != nil {
		fmt.Printf("Failed to subscribe: %v", err)
		return
	}

	for {
		event, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("Stream closed")
				return
			} else {
				log.Fatalf("Error receiving event: %v", err)
			}
		} else {
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.processEvent(event)
			}()
		}
	}
}

func (c *CloudeventsServiceClient) processEvent(event *v1.CloudEvent) {
	c.Lock()
	defer c.Unlock()
	log.Printf("Received event: %+v", event)
	counterEvents++
	log.Println("Отловлено: ", counterEvents)
}
