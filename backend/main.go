package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
	jwtKey = []byte("secret_key") // Use a secure key in production
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint
	Text      string
	CreatedAt time.Time
}

func init() {
	var err error

	dsn := os.Getenv("DB_DSN")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	db.AutoMigrate(&User{}, &Message{})

	rdb = redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}

func main() {
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Register Endpoint"))
	})
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Login Endpoint"))
	})
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
        
