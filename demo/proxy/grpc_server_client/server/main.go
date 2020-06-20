package main

import (
	"context"
	"flag"
	"fmt"
	proto "gateway_demo/demo/proxy/grpc_server_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

var port = flag.Int("port", 50055, "the port to serve on")

const (
	streamingCount = 10
)

type server struct {
}

func (s *server) UnaryEcho(ctx context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {
	fmt.Printf("----- UnaryEcho ----\n")
	metadata.FromIncomingContext(ctx)
}

func (s *server) ServerStreamingEcho(in *proto.EchoRequest, stream proto.Echo_ServerStreamingEchoServer) error {

}

func (s *server) ClientStreamingEcho(stream proto.Echo_ClientStreamingEchoServer) error {

}

func (s *server) BidirectionalStreamEcho(steam proto.Echo_BidirectionalStreamEchoServer) error {

}

func main() {
	flag.Parse()
	s := grpc.NewServer()

	proto.RegisterEchoServer(s, &server{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s.Serve(lis)

}
