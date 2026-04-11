package main

import (
	"context"
	pb "grpc_pingpong/pingpong"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	ping pb.UnimplementedPingPongServiceServer
}

func (s *server) SendPing(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	log.Println("Received message: ", req.GetMessage())
	return &pb.PongResponse{Message: "Pong!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Println("failed to listen: ", err)
		return
	}

	s := grpc.NewServer()
	var NewServer server
	pb.RegisterPingPongServiceServer(s, &NewServer.ping)

	log.Println("Server listening to port: 5001")
	if err := s.Serve(lis); err != nil {
		log.Println("failed to server: ", err)
		return
	}
}
