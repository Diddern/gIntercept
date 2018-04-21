package main

import (
	"log"
	"net"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/Diddern/gIntercept/pb"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
)

type server struct{}

func main()  {
	var portNumber = ":5001"
	lis, err := net.Listen("tcp", portNumber)
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Print("Listening on port " + portNumber)
	s := grpc.NewServer(grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	pb.RegisterGCDServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
func (s *server) Compute(ctx context.Context, r *pb.GCDRequest) (*pb.GCDResponse, error) {
	a, b := r.A, r.B
	for b != 0 {
		a, b = b, a%b
	}
	return &pb.GCDResponse{Result: a}, nil
}