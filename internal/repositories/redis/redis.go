package redis

import (
	"github.com/andReyM228/lib/log"
	"github.com/andReyM228/lib/redis"
)

type Repository struct {
	log   log.Logger
	redis redis.Redis
}

func NewRepository(redis redis.Redis, log log.Logger) Repository {
	return Repository{
		redis: redis,
		log:   log,
	}
}

func (r Repository) Create(key string, value interface{}) error {
	return r.redis.Set(key, value)
}

func (r Repository) GetString(key string) (string, error) {
	return r.redis.GetString(key)
}

func (r Repository) GetBytes(key string) ([]byte, error) {
	return r.redis.GetBytes(key)
}
