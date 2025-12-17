
# MangaHub Architecture
## 1) Overview
MangaHub is a multi-protocol, multi-service system:
- HTTP REST API is the central orchestrator (auth, DB operations, user actions)
- TCP is used for real-time progress synchronization broadcast
- UDP is used for lightweight chapter-release notifications
- WebSocket is used for real-time chat broadcasting
- gRPC is used for internal service calls (GetManga/SearchManga/UpdateProgress)

All services share the same data model (SQLite) to ensure consistent state.

---

## 2) Components & Ports

| Component | Protocol | Port | Responsibility |
|---|---|---:|---|
| api-server | HTTP (Gin) | 8080 | REST API + JWT + DB + triggers TCP/UDP |
| tcp-server | TCP | 9090 | Broadcast progress updates to connected clients |
| udp-server | UDP | 9091 | Register clients + broadcast chapter notificat ions |
| grpc-server | gRPC | 9092 | Internal RPC: GetManga/SearchManga/UpdateProgress |
| websocket-server | WebSocket | 9093 | Chat: broadcast messages to connected users |

---

## 3) Data Layer (SQLite)
Minimum schema:
- users(id, username, email, password_hash, created_at)
- manga(id, title, author, genres, status, total_chapters, description, cover_url)
- user_progress(user_id, manga_id, current_chapter, status, updated_at)

Seeding:
- On API startup, if manga table is empty, load `data/manga.json`.

---

## 4) Integration Flows (Mandatory)

### 4.1 HTTP Progress Update ⇒ TCP Broadcast
1) Client calls HTTP `PUT /users/progress` with JWT.
2) API server updates `user_progress` in SQLite.
3) API server sends a newline-delimited JSON ProgressUpdate to TCP server `127.0.0.1:9090`.
4) TCP server broadcasts the update to all connected TCP clients.

### 4.2 UDP Notifications
1) UDP clients send `{"type":"register","user_id":"..."}` to `:9091`.
2) API triggers notification via HTTP `POST /admin/notify`.
3) API sends a UDP control packet `{"type":"broadcast","manga_id":"...","message":"..."}` to UDP server.
4) UDP server broadcasts `{"type":"chapter_release",...}` to all registered clients.

### 4.3 WebSocket Chat
1) Clients connect to `ws://127.0.0.1:9093/ws`.
2) Any client message is broadcast to all connected clients.

### 4.4 gRPC Internal Services
- `GetManga(id)` reads from manga table.
- `SearchManga(query)` searches by title/author.
- `UpdateProgress(user_id, manga_id, chapter)` writes to `user_progress` and triggers TCP notifier (same as HTTP flow).

---

## 5) Mapping Use Cases to Protocols

| Use Case | Protocol(s) | Implementation |
|---|---|---|
| Register/Login | HTTP + JWT | /auth/register, /auth/login |
| Browse/Search Manga | HTTP / gRPC | GET /manga, GET /manga/:id, gRPC Search/Get |
| Add to Library | HTTP | POST /users/library |
| Update Progress | HTTP / gRPC + TCP | PUT /users/progress or gRPC UpdateProgress ⇒ TCP broadcast |
| Chapter Release Notification | HTTP + UDP | POST /admin/notify ⇒ UDP broadcast |
| Chat | WebSocket | ws://.../ws |

---

## 6) Operational Notes
- Run all 5 servers for a full demo.
- TCP/UDP/WS are designed as dedicated processes to demonstrate protocol usage clearly.
