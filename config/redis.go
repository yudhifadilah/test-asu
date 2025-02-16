package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func ConnectRedis() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Gagal membaca file .env, menggunakan variabel lingkungan")
	}

	redisURL := os.Getenv("REDIS_URL")

	var options *redis.Options
	var err error

	if redisURL != "" {
		// Jika menggunakan URL Redis dengan TLS/SSL (rediss://)
		options, err = redis.ParseURL(redisURL)
		if err != nil {
			return fmt.Errorf("❌ Gagal parsing Redis URL: %w", err)
		}
		options.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Gunakan TLS
	} else {
		// Jika menggunakan konfigurasi manual
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")
		redisPassword := os.Getenv("REDIS_PASSWORD")

		options = &redis.Options{
			Addr:      fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password:  redisPassword,
			DB:        0,
			TLSConfig: &tls.Config{InsecureSkipVerify: true}, // Gunakan TLS
		}
	}

	// Inisialisasi Redis Client
	RedisClient = redis.NewClient(options)

	// Test koneksi Redis
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("❌ Gagal menghubungkan ke Redis: %w", err)
	}

	log.Println("✅ Redis terhubung dengan sukses melalui TLS!")
	return nil
}
