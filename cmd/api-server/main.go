package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"mangahub/internal/auth"
	tcpmod "mangahub/internal/tcp"
	udpmod "mangahub/internal/udp"
	"mangahub/internal/user"
	"mangahub/pkg/database"
	
)

func main() {
	// ===== Config =====
	dbPath := getenv("MANGAHUB_DB", "./mangahub.db")
	jwtSecret := getenv("MANGAHUB_JWT_SECRET", "dev-secret")

	// ===== DB =====
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	userRepo := user.NewRepo(db)

	// ===== External services clients =====
	tcpNotifier := tcpmod.NewNotifier("127.0.0.1:9090")
	udpControl := udpmod.NewControl("127.0.0.1:9091")

	// ===== HTTP Router =====
	r := gin.Default()

	authed := r.Group("/", auth.JWTMiddleware(jwtSecret))

	// ===== PUT /users/progress (HTTP → TCP) =====
	authed.PUT("/users/progress", func(c *gin.Context) {
		var req user.UpdateProgressRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		userID := auth.GetUserID(c)

		if err := userRepo.UpdateProgress(
			userID,
			req.MangaID,
			req.Chapter,
			req.Status,
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Trigger TCP broadcast (bắt buộc theo spec)
		_ = tcpNotifier.SendProgressUpdate(tcpmod.ProgressUpdate{
			UserID:    userID,
			MangaID:   req.MangaID,
			Chapter:   req.Chapter,
			Timestamp: time.Now().Unix(),
		})

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// ===== POST /admin/notify (HTTP → UDP) =====
	authed.POST("/admin/notify", func(c *gin.Context) {
		var req udpmod.NotifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		if req.MangaID == "" || req.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing manga_id or message"})
			return
		}
		if err := udpControl.Broadcast(req.MangaID, req.Message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	log.Println("HTTP API listening on :8080")
	log.Fatal(r.Run(":8080"))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
