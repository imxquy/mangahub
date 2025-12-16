package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "mangahub/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMangaServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1) GetManga
	fmt.Println("=== GetManga(one-piece) ===")
	gm, err := client.GetManga(ctx, &pb.GetMangaRequest{Id: "one-piece"})
	if err != nil {
		log.Fatal("GetManga:", err)
	}
	fmt.Println(gm)

	// 2) Search
	fmt.Println("=== Search(query=piece) ===")
	sr, err := client.SearchManga(ctx, &pb.SearchRequest{Query: "piece", Limit: 5})
	if err != nil {
		log.Fatal("Search:", err)
	}
	fmt.Println(sr)

	// 3) UpdateProgress (phải trigger TCP broadcast nếu bạn đã làm đúng server handlers)
	fmt.Println("=== UpdateProgress(user=u1, manga=one-piece, ch=2) ===")
	up, err := client.UpdateProgress(ctx, &pb.ProgressRequest{UserId: "u1", MangaId: "one-piece", Chapter: 2})
	if err != nil {
		log.Fatal("UpdateProgress:", err)
	}
	fmt.Println(up)
}
