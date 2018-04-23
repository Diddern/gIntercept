package main

import (
	"log"
	"net"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/Diddern/gIntercept/pb"
	"google.golang.org/grpc/credentials"

)

type server struct{}

func main()  {

	portNumber := ":5001"
	creds, err := credentials.NewServerTLSFromFile("../gRPC-simpleGCDService/gcd/server-cert.pem", "../gRPC-simpleGCDService/gcd/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}
	lis, err := net.Listen("tcp", portNumber)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Print("Listening on port 5001:")
	s := grpc.NewServer(
		grpc.Creds(creds),)
	pb.RegisterGCDServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (s *server) Compute(ctx context.Context, requestFromClient *pb.GCDRequest) (*pb.GCDResponse, error) {

	address := "localhost:3000"
	creds, err := credentials.NewClientTLSFromFile("../gRPC-simpleGCDService/gcd/server-cert.pem", "")
	if err != nil {
		log.Fatalf("cert load error: %s", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil{
		log.Fatalf("Dail failed %v", err)
	}
	log.Print("Dailed sucessfully to ", address)

	gcdClient := pb.NewGCDServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resultFromServer, err := gcdClient.Compute(ctx, requestFromClient,)
	if err != nil{
		log.Fatalf("Could not send to server: ", err)
	}
	log.Print("ctx: ", requestFromClient)
	log.Print("Request from client: ", requestFromClient)
	log.Print("Response from server: ", resultFromServer)
	return &pb.GCDResponse{Result: resultFromServer.Result}, nil
}

