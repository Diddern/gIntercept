package main

import (
	"log"
	"google.golang.org/grpc"
	"github.com/Diddern/gIntercept/pb"
	"net"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc/reflection"
)

func main()  {
	var portNumber = ":5001"
	lis, err := net.Listen("tcp", portNumber)
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Print("Listening on port " + portNumber)
	s := grpc.NewServer()
	pb.RegisterdecoderServiceServer(s, &server{})
	reflection.Register(s)

}