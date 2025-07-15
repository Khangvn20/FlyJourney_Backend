package service

import (
	"time"
)
type RedisService interface {
    Set(key string, value interface{}, expiration time.Duration) error
    Get(key string) (string, error)
    Del(key string) error
    Exists(key string) (bool, error)
}