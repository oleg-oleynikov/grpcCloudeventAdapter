package main

import (
	"context"
	"fmt"

	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	protoCloudevents "open-cluster-management.io/sdk-go/pkg/cloudevents/generic/options/grpc/protobuf/v1"
)

type CloudEventServiceServer struct {
	sync.RWMutex
	protoCloudevents.UnimplementedCloudEventServiceServer

	server      *grpc.Server
	subscribers map[protoCloudevents.CloudEventService_SubscribeServer]chan *protoCloudevents.CloudEvent
	eventChan   chan *protoCloudevents.CloudEvent
	wg          sync.WaitGroup
}

func NewCloudEventServiceServer(opt ...grpc.ServerOption) *CloudEventServiceServer {
	return &CloudEventServiceServer{
		server:      grpc.NewServer(opt...),
		subscribers: make(map[protoCloudevents.CloudEventService_SubscribeServer]chan *protoCloudevents.CloudEvent),
		eventChan:   make(chan *protoCloudevents.CloudEvent, 100),
	}
}

func (s *CloudEventServiceServer) Publish(ctx context.Context, req *protoCloudevents.PublishRequest) (*emptypb.Empty, error) {
	event := req.GetEvent()
	log.Println("Пришло нахуй")
	log.Println(event)
	select {
	case s.eventChan <- event:
	default:
		log.Printf("eventChan is full, dropping event: %v", event)
		return nil, fmt.Errorf("eventChan is full, event dropped")
	}

	go s.processEvents()
	return &empty.Empty{}, nil
}

func (s *CloudEventServiceServer) processEvents() {
	for {
		select {
		case event := <-s.eventChan:
			s.RLock()
			s.wg.Add(1)
			go func(event *protoCloudevents.CloudEvent) {
				defer s.wg.Done()
				s.processEvent(event)
			}(event)
			s.RUnlock()
		default:
			return
		}
	}
}

func (s *CloudEventServiceServer) processEvent(event *protoCloudevents.CloudEvent) {
	for _, eventChan := range s.subscribers {
		select {
		case eventChan <- event:

		default:
			log.Printf("Subscriber queue is full, dropping event")
		}
	}
}

func (s *CloudEventServiceServer) addSubscriber(stream protoCloudevents.CloudEventService_SubscribeServer) {
	s.Lock()
	defer s.Unlock()

	eventChan := make(chan *protoCloudevents.CloudEvent, 100)
	s.subscribers[stream] = eventChan

	go func() {
		for {
			event, ok := <-eventChan
			if !ok {
				return
			}
			if err := stream.Send(event); err != nil {
				log.Printf("Failed to send event to subscriber: %v", err)
				s.removeSubscriber(stream)
				return
			}
		}
	}()
}

func (s *CloudEventServiceServer) removeSubscriber(subscriber protoCloudevents.CloudEventService_SubscribeServer) {
	s.Lock()
	defer s.Unlock()

	if eventChan, ok := s.subscribers[subscriber]; ok {
		subscriber.Context().Done()
		close(eventChan)
		delete(s.subscribers, subscriber)
	}
}

func (s *CloudEventServiceServer) Subscribe(req *protoCloudevents.SubscriptionRequest, stream protoCloudevents.CloudEventService_SubscribeServer) error {
	log.Printf("New subscriber: %v", req)
	s.addSubscriber(stream)

	<-stream.Context().Done()

	err := stream.Context().Err()
	if err != nil {
		log.Printf("Subscriber disconnected: %v\n", err)
	}

	s.removeSubscriber(stream)
	return nil
}

func (s *CloudEventServiceServer) Run(network string, addr string) {
	lis, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	protoCloudevents.RegisterCloudEventServiceServer(s.server, s)

	log.Printf("Starting listening server on %s\n", addr)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			fmt.Printf("Failed to serve: %v", err)
		}
	}()
}

func (s *CloudEventServiceServer) Shutdown() {
	s.server.GracefulStop()

	s.RLock()
	for sub := range s.subscribers {
		s.removeSubscriber(sub)
	}
	s.RUnlock()
	s.wg.Wait()
}

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// corsHandler := cors.NewCorsGrpcBuilder().WithAllowedOrigins("localhost:3000").BuildHandler()
	// unaryInterceptor := cors.BuildCorsUnaryInterceptor(corsHandler)
	// streamInterceptor := cors.BuildCorsStreamInterceptor(corsHandler)
	// c := NewCloudEventServiceServer(unaryInterceptor, streamInterceptor)
	c := NewCloudEventServiceServer(grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// log.Println("Типа перехватил")
		return handler(ctx, req)
	}))

	c.Run("tcp", ":50051")

	<-quit
	c.Shutdown()
}
