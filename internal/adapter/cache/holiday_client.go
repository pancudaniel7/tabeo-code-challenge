package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"tabeo.org/challenge/internal/core/entity"
	"tabeo.org/challenge/internal/pkg/apperr"
)

type RedisHolidayCacheClient struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewHolidayCacheClient() HolidayCacheClient {
	host := viper.GetString("cache.host")
	port := viper.GetInt("cache.port")
	db := viper.GetInt("cache.db")
	password := viper.GetString("cache.password")
	ttl := viper.GetInt("cache.ttl")
	if ttl == 0 {
		ttl = 86400
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	return &RedisHolidayCacheClient{
		rdb: rdb,
		ttl: time.Duration(ttl) * time.Second,
	}
}

func (c *RedisHolidayCacheClient) key(year int, country string) string {
	return fmt.Sprintf("publicholidays:%d:%s", year, country)
}

func (c *RedisHolidayCacheClient) GetPublicHolidays(ctx context.Context, year int, country string) ([]entity.PublicHolidays, error) {
	val, err := c.rdb.Get(ctx, c.key(year, country)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, apperr.NotFoundErr("cache miss", nil)
	} else if err != nil {
		return nil, apperr.Internal("redis get failed", err)
	}

	var resp []PublicHolidaysCacheDTO
	err = json.Unmarshal([]byte(val), &resp)
	if err != nil {
		return nil, apperr.Internal("failed to unmarshal cached holidays", err)
	}

	holidays := make([]entity.PublicHolidays, len(resp))
	for i, h := range resp {
		holidays[i] = h.ToEntity()
	}
	return holidays, nil
}

func (c *RedisHolidayCacheClient) SetPublicHolidays(ctx context.Context, year int, country string, holidays []entity.PublicHolidays) error {
	resp := make([]*PublicHolidaysCacheDTO, len(holidays))
	for i, h := range holidays {
		resp[i] = resp[i].ToDTO(&h)
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return apperr.Internal("failed to marshal holidays for cache", err)
	}

	status := c.rdb.Set(ctx, c.key(year, country), data, c.ttl)
	if status.Err() != nil {
		return apperr.Internal("redis set failed", status.Err())
	}
	return nil
}
