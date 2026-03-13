package redis

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
	ctx    context.Context
}

func NewStorage(addr string, password string, db int) (*Storage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Storage{
		client: rdb,
		ctx:    context.Background(),
	}, nil
}

func (r *Storage) Save(stat models.Stat) error {
	data, err := json.Marshal(stat)
	if err != nil {
		return fmt.Errorf("failed to marshal data stat: %w", err)
	}

	key := fmt.Sprintf("stat:%s%s %d", stat.Base, stat.Quote, stat.Timedump)

	err = r.client.Set(r.ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to save to redis: %w", err)
	}

	return nil
}

func (r *Storage) GetStat() ([]models.Stat, error) {
	// 1. Получаем все ключи по маске
	// В Redis лучше использовать Scan, чтобы не блокировать базу
	ctx := context.Background()
	iter := r.client.Scan(ctx, 0, "stat:*", 0).Iterator()

	var stats []models.Stat
	for iter.Next(ctx) {
		key := iter.Val()

		// 2. Достаем значение по ключу
		val, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Пропускаем, если ключ исчез (например, истек TTL)
		}

		// 3. Декодируем JSON (в Redis мы обычно храним структуру как JSON-строку)
		var s models.Stat
		if err := json.Unmarshal([]byte(val), &s); err != nil {
			continue
		}

		stats = append(stats, s)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	// 4. Сортировка (так как Redis выдает ключи в случайном порядке)
	// Аналог SQL-ного "ORDER BY timedump DESC"
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Timedump.After(stats[j].Timedump)
	})

	// 5. Лимит (аналог LIMIT 100)
	if len(stats) > 100 {
		stats = stats[:100]
	}

	return stats, nil
}
