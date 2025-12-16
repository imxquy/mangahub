package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	pb "mangahub/proto"
	tcpmod "mangahub/internal/tcp"
	"mangahub/pkg/database"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMangaServiceServer
	db          *sql.DB
	tcpNotifier *tcpmod.Notifier
}

func main() {
	// ===== Config =====
	dbPath := getenv("MANGAHUB_DB", "./mangahub.db")

	// ===== DB =====
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// ===== TCP notifier (same as HTTP flow) =====
	tcpNotifier := tcpmod.NewNotifier("127.0.0.1:9090")

	// ===== gRPC listener =====
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer()

	// ===== Register gRPC service with injected deps =====
	pb.RegisterMangaServiceServer(s, &server{
		db:          db,
		tcpNotifier: tcpNotifier,
	})

	log.Println("gRPC listening on :9092")
	log.Fatal(s.Serve(lis))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
