package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ─── Puerto de Caché ────────────────────────────────────────────────────────

// CachePort define la interfaz de caché que consumirán los casos de uso.
// Si la implementación de Redis falla, los consumidores deben hacer fallback a BD.
type CachePort interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Invalidate(ctx context.Context, key string) error
	InvalidatePrefix(ctx context.Context, prefix string) error
}

// ─── Implementación Redis ───────────────────────────────────────────────────

// RedisCache implementa CachePort usando go-redis/v9.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache crea una nueva instancia conectada a Redis.
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error al conectar con Redis: %w", err)
	}

	fmt.Println("Conexion a Redis exitosa")
	return &RedisCache{client: client}, nil
}

// Close cierra la conexión con Redis.
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// Set serializa el valor a JSON y lo almacena con un TTL.
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache set: error al serializar: %w", err)
	}
	return rc.client.Set(ctx, key, data, ttl).Err()
}

// Get obtiene el valor de Redis y lo deserializa en dest.
// Retorna redis.Nil si la key no existe (cache miss).
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		return err // redis.Nil para cache miss, otro error para fallo real
	}
	return json.Unmarshal(data, dest)
}

// Invalidate elimina una key específica.
func (rc *RedisCache) Invalidate(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

// InvalidatePrefix elimina todas las keys que coincidan con un patrón glob.
// Ejemplo: "acts:search:*" borra todas las búsquedas cacheadas de actos.
func (rc *RedisCache) InvalidatePrefix(ctx context.Context, prefix string) error {
	iter := rc.client.Scan(ctx, 0, prefix, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return rc.client.Del(ctx, keys...).Err()
	}
	return nil
}
