# MangaHub — Manga & Comic Tracking System (Net Centric Programming)

This repository implements the MangaHub term project end-to-end using:
- HTTP REST API + JWT + SQLite (port **8080**)
- TCP Progress Sync Broadcast (port **9090**)
- UDP Chapter Release Notifications (port **9091**)
- gRPC Internal Services (port **9092**)
- WebSocket Chat (port **9093**)

Ports follow the CLI manual defaults: 8080/9090/9091/9092/9093.

---

## 1) Tech Stack
- Go
- Gin (HTTP REST)
- SQLite (mattn/go-sqlite3)
- JWT (golang-jwt/jwt)
- TCP sockets (net)
- UDP sockets (net)
- WebSocket (gorilla/websocket)
- gRPC (google.golang.org/grpc + protobuf)

---

## 2) Repository Structure

cmd/
api-server/ # HTTP REST + JWT + SQLite + TCP/UDP integration
tcp-server/ # TCP progress broadcast service
udp-server/ # UDP register + broadcast service
grpc-server/ # gRPC internal services
websocket-server/ # WebSocket chat service

tcp-monitor/ # (optional) dev client to monitor TCP broadcasts
udp-client/ # (optional) dev client to receive UDP notifications
ws-client/ # (optional) dev client to send/receive WS messages
grpc-client/ # (optional) dev client to call gRPC methods

internal/
auth/ # JWT middleware + auth service/handlers
manga/ # manga repo/handlers + seeding
user/ # library/progress repo/handlers
tcp/ # TCP server + notifier client
udp/ # UDP server + control client
websocket/ # hub/client/message
grpc/ # (if split) grpc helpers

pkg/
database/ # SQLite Open/Migrate
models/ # shared models (optional)

proto/
manga.proto # gRPC service definitions

data/
manga.json # seed data


---

## 3) Prerequisites (Windows)
- Go installed
- `protoc` installed (only needed if regenerating proto)

Install deps:
```powershell
go mod tidy
