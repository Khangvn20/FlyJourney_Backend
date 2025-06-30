package config

import (
    "os"
    "strconv"
    "github.com/redis/go-redis/v9"
    "fmt"
    "context"
    "time"
)

type RedisConfig struct {
    Host     string
    Port     string
    Password string
    DB       int
}

func NewRedisConfig() *RedisConfig {
    dbStr := os.Getenv("REDIS_DB")
    db, err := strconv.Atoi(dbStr)
    if err != nil {
        db = 0 
    }

    return &RedisConfig{
        Host:     os.Getenv("REDIS_HOST"),
        Port:     os.Getenv("REDIS_PORT"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       db,
    }
}

func NewRedisClient(cfg *RedisConfig) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
    })
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := client.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to connect %w", err)
    }

    return client, nil
}