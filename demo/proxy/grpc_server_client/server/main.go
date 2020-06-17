package main

import (
	"context"
	"flag"
	proto "gateway_demo/demo/proxy/grpc_server_client"
)

var port = flag.Int("port", 50055, "the port to serve on")

const (
	streamingCount = 10
)

type server struct {
}

func (s *server) UnaryEcho(ctx context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {

}

func (s *server) ServerStreamingEcho(in *proto.EchoRequest, stream proto.Echo_ServerStreamingEchoServer) error {

}

func (s *server) ClientStreamingEcho(stream proto.Echo_ServerStreamingEchoClient) error {

}

func (s *server) BidirectionalStreamEcho(steam proto.Echo_ServerStreamingEchoServer) error {

}

func main() {

}
