# MangaHub Demo Checklist (5 Protocols)

## Pre-demo setup
- [ ] Start TCP server (:9090)
- [ ] Start UDP server (:9091)
- [ ] Start WebSocket server (:9093)
- [ ] Start gRPC server (:9092)
- [ ] Start HTTP API server (:8080)
- [ ] Start TCP monitor client (optional) to display broadcasts
- [ ] Start UDP client (optional) to display notifications
- [ ] Start 2 WebSocket clients (optional) for chat demo

---

## 1) HTTP + JWT + SQLite (Core)
- [ ] Register: POST /auth/register
- [ ] Login: POST /auth/login (obtain JWT)
- [ ] Manga list/search: GET /manga?q=...
- [ ] Manga details: GET /manga/:id
- [ ] Library add: POST /users/library (JWT)
- [ ] Library list: GET /users/library (JWT)

---

## 2) Mandatory Integration: HTTP Progress ⇒ TCP Broadcast
- [ ] Run tcp-monitor connected to :9090
- [ ] Call PUT /users/progress (JWT)
- [ ] tcp-monitor prints broadcast JSON

---

## 3) UDP Notifications
- [ ] Run udp-client and send register packet
- [ ] Call POST /admin/notify (JWT)
- [ ] udp-client receives chapter_release JSON

---

## 4) WebSocket Chat
- [ ] Open 2 ws clients to ws://127.0.0.1:9093/ws
- [ ] Send message from client A
- [ ] Client B receives broadcast message

---

## 5) gRPC Internal Services
- [ ] GetManga(id) returns record
- [ ] SearchManga(query) returns results
- [ ] UpdateProgress triggers TCP broadcast (observe tcp-monitor)

---

## End
- [ ] All five protocols demonstrated successfully on required ports.
