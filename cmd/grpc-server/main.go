package main

import (
	"log"
	"net"

	pb "mangahub/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMangaServiceServer
	// TODO: inject db + tcp broadcaster later
}

func main() {
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMangaServiceServer(s, &server{})

	log.Println("gRPC listening on :9092")
	log.Fatal(s.Serve(lis))
}
