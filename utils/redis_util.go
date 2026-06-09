package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"nav-rtlogging-go/global"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	client *redis.Client
	ctx    context.Context
}

var (
	RedisUtilApp           *RedisUtil
	once                   sync.Once
	REDIS_RTLOGGING_PREFIX = "radar_rtlogging_"
	RedisNotExpire         = time.Duration(0) * time.Second
	RedisDayExpire         = time.Duration(24*60*60) * time.Second
	RedisHourExpire        = time.Duration(60*60) * time.Second
	RedisOpTimeout         = 3 * time.Second
)

func InitRedis() *RedisUtil {
	redisConfig := global.NAV_CONFIG.Redis
	redisUtil := GetInstance(fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port), redisConfig.Password, redisConfig.DB)
	return redisUtil
}

func GetInstance(addr, password string, db int) *RedisUtil {
	once.Do(func() {
		RedisUtilApp = &RedisUtil{
			client: redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
				DB:       db,
			}),
			ctx: context.Background(),
		}
	})
	return RedisUtilApp
}

// Set 设置值（带过期时间）
func (r *RedisUtil) Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	if expiration <= 0 {
		return r.client.Set(ctx, key, value, 0).Err()
	}
	return r.client.Set(ctx, key, value, expiration).Err()
}

// SetForever 设置永久不过期的值
func (r *RedisUtil) SetForever(key string, value interface{}) error {
	return r.Set(key, value, 0)
}

// GetObj 泛型 Get 方法
func (r *RedisUtil) GetObj(key string, dest interface{}) error {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return err
	}
	return nil
}

// SetObj 对应的 Set 方法，需要存储为 JSON
func (r *RedisUtil) SetObj(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	return r.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取值
func (r *RedisUtil) Get(key string) (string, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	return r.client.Get(ctx, key).Result()
}

// Del 删除 key
func (r *RedisUtil) Del(keys ...string) error {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	return r.client.Del(ctx, keys...).Err()
}

// Exists 判断 key 是否存在
func (r *RedisUtil) Exists(key string) (bool, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Expire 设置 key 的过期时间
func (r *RedisUtil) Expire(key string, expiration time.Duration) (bool, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	return r.client.Expire(ctx, key, expiration).Result()
}

// Incr 自增 key
func (r *RedisUtil) Incr(key string) (int64, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	return r.client.Incr(ctx, key).Result()
}

// Decr 自减 key
func (r *RedisUtil) Decr(key string) (int64, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()
	return r.client.Decr(ctx, key).Result()
}

func (r *RedisUtil) contextWithTimeout() (context.Context, context.CancelFunc) {
	if RedisOpTimeout <= 0 {
		return context.WithCancel(r.ctx)
	}
	return context.WithTimeout(r.ctx, RedisOpTimeout)
}
