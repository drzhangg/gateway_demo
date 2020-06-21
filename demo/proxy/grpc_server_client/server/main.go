package main

import (
	"context"
	"flag"
	"fmt"
	proto "gateway_demo/demo/proxy/grpc_server_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}
	fmt.Println("md", md)
	fmt.Printf("request received: %v, sending echo\n", in)
	return &proto.EchoResponse{Message: in.Message}, nil
}

func (s *server) ServerStreamingEcho(in *proto.EchoRequest, stream proto.Echo_ServerStreamingEchoServer) error {
	fmt.Printf("---- ServerStreamingEcho ---\n")
	fmt.Printf("request received: %v\n", in)
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", in.Message)
		err := stream.Send(&proto.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ClientStreamingEcho(stream proto.Echo_ClientStreamingEchoServer) error {
	fmt.Printf("--- ClientStreamingEcho ---\n")

	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&proto.EchoResponse{Message: message})
		}
		message = in.Message
		fmt.Printf("request received：%v, building echo\n", in)
		if err != nil {
			return err
		}
	}
}

func (s *server) BidirectionalStreamEcho(stream proto.Echo_BidirectionalStreamEchoServer) error {
	fmt.Printf("--- BidirectionalStreamingEcho ---\n")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		if err := stream.Send(&proto.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
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
