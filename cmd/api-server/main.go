package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"mangahub/internal/auth"
	mangamod "mangahub/internal/manga"
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

	// ===== Seed manga if empty =====
	if err := mangamod.SeedIfEmpty(db, "./data/manga.json"); err != nil {
		log.Fatalf("seed manga: %v", err)
	}

	// ===== Services & repos =====
	authSvc := auth.NewService(db, jwtSecret)
	mangaRepo := mangamod.NewRepo(db)
	userRepo := user.NewRepo(db)

	// ===== External services clients =====
	tcpNotifier := tcpmod.NewNotifier("127.0.0.1:9090")
	udpControl := udpmod.NewControl("127.0.0.1:9091")

	// ===== HTTP Router =====
	r := gin.Default()

	// ===== Public routes =====
	auth.RegisterRoutes(r, authSvc)
	mangamod.RegisterRoutes(r, mangaRepo)

	// ===== Protected routes =====
	authed := r.Group("/", auth.JWTMiddleware(jwtSecret))

	// Library routes
	user.RegisterLibraryRoutes(authed, userRepo)

	// PUT /users/progress (HTTP → DB → TCP broadcast)
	authed.PUT("/users/progress", func(c *gin.Context) {
		var req user.UpdateProgressRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		userID := auth.GetUserID(c)
		if err := userRepo.UpdateProgress(userID, req.MangaID, req.Chapter, req.Status); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// trigger TCP broadcast (best-effort)
		_ = tcpNotifier.SendProgressUpdate(tcpmod.ProgressUpdate{
			UserID:    userID,
			MangaID:   req.MangaID,
			Chapter:   req.Chapter,
			Timestamp: time.Now().Unix(),
		})

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// POST /admin/notify (HTTP → UDP broadcast)
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
