package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mouradev1/buscacepsgolang/internal/config"
	"github.com/mouradev1/buscacepsgolang/internal/models"
)

func GetCepCache(ctx context.Context, cep string) (*models.Cep, error) {
	key := fmt.Sprintf("cep:%s", cep)
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if val == "null" {
		return nil, nil
	}
	var cepData models.Cep
	if err := json.Unmarshal([]byte(val), &cepData); err != nil {
		return nil, err
	}
	return &cepData, nil
}

func SetCepCache(ctx context.Context, cep string, data *models.Cep, ttl time.Duration) error {
	key := fmt.Sprintf("cep:%s", cep)
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return config.RedisClient.Set(ctx, key, bytes, ttl).Err()
}

func SetCepNotFoundCache(ctx context.Context, cep string, ttl time.Duration) error {
	key := fmt.Sprintf("cep:%s", cep)
	return config.RedisClient.Set(ctx, key, "null", ttl).Err()
}
