package grpcutil

import (
	"context"
	"log"
	"path"
	"time"

	"google.golang.org/grpc"
)

// RequestLogger ...
func RequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	service := path.Dir(info.FullMethod)[1:]
	method := path.Base(info.FullMethod)
	log.Printf("%s %s %v %s %s\n", service, method, grpc.Code(err), time.Since(start), grpc.ErrorDesc(err))
	return h, err
}
