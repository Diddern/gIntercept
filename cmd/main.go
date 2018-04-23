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
	server := grpc.NewServer()
	pb.RegisterGCDServiceServer(server, &server{})
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
func (s *server) Compute(ctx context.Context, requestFromClient *pb.GCDRequest) (*pb.GCDResponse, error) {

	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil{
		log.Fatalf("Dail failed %v", err)
	}

	gcdClient := pb.NewGCDServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resultFromServer, err := gcdClient.Compute(ctx, requestFromClient,)
	if err != nil{
		log.Fatalf("Could not send to server: %v", err)
	}
	requestFromClient.ProtoMessage()
	log.Print("ctx: ", requestFromClient)
	log.Print("Request from client: ", requestFromClient)
	log.Print("Response from server: ", resultFromServer)


	return &pb.GCDResponse{Result: resultFromServer.Result}, nil
}