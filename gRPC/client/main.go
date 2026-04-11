package main

import (
	"context"
	pb "grpc_pingpong/pingpong"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error in connecting: ", err)
		return
	}

	defer conn.Close()

	c := pb.NewPingPongServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),3 *	 time.Second)
	defer cancel()

	r, err := c.Sending(ctx, &pb.PingRequest{Message: "Ping_to_Server"})
	if err != nil {
		log.Println("could not ping: ", err)
		return
	}

	log.Println("Response from server: ", r.GetMessage())
}
