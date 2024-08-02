package grpccors

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type corsCheck func(md metadata.MD) error
type CorsHandler func(ctx context.Context, method string) error

type CorsGrpc struct {
	corsChecks []corsCheck
}

func (c *CorsGrpc) addCorsCheck(handler corsCheck) {
	c.corsChecks = append(c.corsChecks, handler)
}

func NewCorsGrpcBuilder() *CorsGrpc {
	return &CorsGrpc{}
}

func (c *CorsGrpc) WithAllowedOrigins(origins ...string) *CorsGrpc {
	if len(origins) == 0 {
		log.Println("Allowed origins is empty")
		return nil
	}

	corsCheckAllowedOrigin := func(md metadata.MD) error {
		origin := md.Get("Origin")
		if len(origin) == 0 {
			return status.Errorf(codes.InvalidArgument, "missing 'Origin' header")
		}

		if !containsString(origin[0], origins) {
			return status.Errorf(codes.PermissionDenied, "origin not allowed: %s", origin[0])
		}

		return nil
	}

	c.addCorsCheck(corsCheckAllowedOrigin)
	return c
}

func (c *CorsGrpc) WithAllowedHeaders(headers ...string) *CorsGrpc {
	if len(headers) == 0 {
		log.Printf("allowed headers is empty")
		return c
	}

	return c
}

func (c *CorsGrpc) BuildHandler() CorsHandler {
	return func(ctx context.Context, method string) error {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Internal, "failed to get metadata")
		}

		for _, h := range c.corsChecks {
			if err := h(md); err != nil {
				return err
			}
		}

		return nil
	}
}

func CorsToServerOptions(corsHandler CorsHandler) []grpc.ServerOption {
	unaryCorsInterceptor := BuildCorsUnaryInterceptor(corsHandler)
	streamCorsInterceptor := BuildCorsStreamInterceptor(corsHandler)
	return []grpc.ServerOption{unaryCorsInterceptor, streamCorsInterceptor}
}

func BuildCorsUnaryInterceptor(cors CorsHandler) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := cors(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	})
}

func BuildCorsStreamInterceptor(cors func(ctx context.Context, method string) error) grpc.ServerOption {
	return grpc.StreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := cors(ss.Context(), info.FullMethod); err != nil {
			return err
		}
		return handler(srv, ss)
	})
}

func containsString(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle || s == "*" {
			return true
		}
	}
	return false
}
