package main

import (
	"context"
	"influencer-golang/config"
	"influencer-golang/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Connect to Redis
	if err := config.ConnectRedis(); err != nil {
		log.Fatalf("❌ Gagal terhubung ke Redis: %v", err)
	}
	log.Println("✅ Redis terhubung dengan sukses!")

	// Cek koneksi Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := config.RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Redis tidak merespons: %v", err)
	}
	log.Println("✅ Redis siap digunakan!")

	// Setup Router dengan CORS
	r := gin.Default()

	// Konfigurasi CORS agar lebih fleksibel
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Ubah jika perlu batasan
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Setup Routes
	routes.SetupRoutes(r)

	// Ambil port dari environment atau gunakan default 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Jalankan server
	log.Printf("🚀 Server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
