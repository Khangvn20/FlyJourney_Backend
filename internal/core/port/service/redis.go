package service

import (
	"time"
)
type RedisService interface {
    Set(key string, value interface{}, expiration time.Duration) error
    Get(key string) (string, error)
    Del(key string) error
    Exists(key string) (bool, error)
    TryLock(key string, value string, expiration time.Duration) (bool, error)
    ReleaseLock(key string,value string) error
    Incr(key string) (int64, error)
    Expire(key string, expiration time.Duration) error
}