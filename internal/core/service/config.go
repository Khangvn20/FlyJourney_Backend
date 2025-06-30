package service

import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
    "time"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)

type redisService struct {
    client *redis.Client
}

func NewRedisService(client *redis.Client) service.RedisService {
    return &redisService{
        client: client,
    }
}


func (s *redisService) Set(key string, value interface{}, expiration time.Duration) error {
    ctx := context.Background()

    var strValue string
    switch v := value.(type) {
    case string:
        strValue = v
    default:
        jsonBytes, err := json.Marshal(value)
        if err != nil {
            return err
        }
        strValue = string(jsonBytes)
    }
    
    return s.client.Set(ctx, key, strValue, expiration).Err()
}

func (s *redisService) Get(key string) (string, error) {
    ctx := context.Background()
    return s.client.Get(ctx, key).Result()
}

func (s *redisService) Del(key string) error {
    ctx := context.Background()
    return s.client.Del(ctx, key).Err()
}

func (s *redisService) Exists(key string) (bool, error) {
    ctx := context.Background()
    val, err := s.client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return val > 0, nil
}