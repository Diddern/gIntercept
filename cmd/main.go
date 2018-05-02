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

var portNumberIn = ":5001"
var addressAndPortNumberOut =  "localhost:3000"
var pathToCert = "../gRPC-simpleGCDService/certs/server-cert.pem"
var pathToKey = "../gRPC-simpleGCDService/certs/server-key.pem"

func main()  {

	//Load cert and key from file
	creds, err := credentials.NewServerTLSFromFile(pathToCert, pathToKey)
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}
	//Listen for incoming connections.
	lis, err := net.Listen("tcp", portNumberIn)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//Create gRPC Server
	s := grpc.NewServer(
		grpc.Creds(creds),)
	pb.RegisterServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) Compute(ctx context.Context, requestFromClient *pb.Request) (*pb.Response, error) {


	creds, err := credentials.NewClientTLSFromFile(pathToCert, "")
	if err != nil {
		log.Fatalf("Could not load certificate from file %v", err)
	}

	// Connect securely to GCD service
	conn, err := grpc.Dial(addressAndPortNumberOut, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}
	defer conn.Close()


	gcdClient := pb.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resultFromServer, err := gcdClient.Compute(ctx, requestFromClient)
	if err != nil{
		log.Fatalf("Could not send to server: %v", err)
	}

	LogRequestsAndResponses(ctx, requestFromClient, resultFromServer)
	return &pb.Response{Result: resultFromServer.Result}, nil
}

func LogRequestsAndResponses(ctx context.Context, requestFromClient *pb.Request, resultFromServer *pb.Response){
	log.Printf("Context for reqest: \t%v", ctx)
	log.Printf("Request from client: \t%v", requestFromClient)
	log.Printf("Response from server: \t%v", resultFromServer)
}