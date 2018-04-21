package main

import (
	"log"
	"net"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/Diddern/gIntercept/pb"
	"time"
)

type server struct{}

func main()  {
	var portNumber = ":5001"
	lis, err := net.Listen("tcp", portNumber)
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Print("Listening on port " + portNumber)
	s := grpc.NewServer()
	pb.RegisterGCDServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
func (s *server) Compute(ctx context.Context, r *pb.GCDRequest) (*pb.GCDResponse, error) {

	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil{
		log.Fatalf("Dail failed %v", err)
	}

	gcdClient := pb.NewGCDServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := gcdClient.Compute(ctx, r,)
	if err != nil{
		log.Fatalf("Could not send to server: ", err)
	}
	log.Print("Mottok A fra klient: ", r.A)
	log.Print("Mottok B fra klient: ", r.B)
	log.Print("Mottok svar fra server: ", res.Result)

	return &pb.GCDResponse{Result: res.Result}, nil
}