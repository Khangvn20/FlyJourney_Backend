package service

import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
    "time"
    "errors"
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
func (s *redisService) TryLock(key string, value string, expiration time.Duration) (bool, error) {
    ctx := context.Background()
    success, err := s.client.SetNX(ctx, key, value, expiration).Result()
    if err != nil {
        return false, err
    }
    return success, nil
}
func (s *redisService) ReleaseLock(key string, value string) error {
    ctx := context.Background()
    
  
    script := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `
    
    result, err := s.client.Eval(ctx, script, []string{key}, value).Result()
    if err != nil {
        return err
    }
    if result.(int64) == 0 {
        return errors.New("lock not found or not owned by caller")
    }
    
    return nil
}
func (s *redisService) Incr(key string) (int64, error) {
    ctx := context.Background()
    return s.client.Incr(ctx, key).Result()
}

func (s *redisService) Expire(key string, expiration time.Duration) error {
    ctx := context.Background()
    success, err := s.client.Expire(ctx, key, expiration).Result()
    if err != nil {
        return err
    }
    if !success {
        return errors.New("failed to set expiration for key")
    }
    return nil
}
func (r *redisService) SetJSON(key string, value interface{}, expiration time.Duration) error {
    ctx := context.Background()
    jsonData, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, jsonData, expiration).Err()
}
func (r *redisService) Keys(pattern string) ([]string, error) {
    ctx := context.Background()
    keys, err := r.client.Keys(ctx, pattern).Result()
    if err != nil {
        return nil, err
    }
    return keys, nil
}