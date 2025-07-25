package redis

import (
	"time"

	"github.com/Ablebil/eco-sample/config"
	"github.com/gofiber/storage/redis"
)

type RedisItf interface {
	SetOTP(email string, otp string, exp time.Duration) error
	GetOTP(email string) (string, error)
	DeleteOTP(email string) error
	SetOAuthState(state string, value []byte, exp time.Duration) error
	GetOAuthState(state string) ([]byte, error)
	DeleteOAuthState(state string) error
}

type Redis struct {
	store *redis.Storage
}

func NewRedis(cfg *config.Config) RedisItf {
	return &Redis{
		store: redis.New(redis.Config{
			Host:     cfg.RedisHost,
			Port:     cfg.RedisPort,
			Password: cfg.RedisPassword,
		}),
	}
}

func (r *Redis) SetOTP(email string, otp string, exp time.Duration) error {
	key := "otp:" + email
	return r.store.Set(key, []byte(otp), exp)
}

func (r *Redis) GetOTP(email string) (string, error) {
	key := "otp:" + email
	val, err := r.store.Get(key)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (r *Redis) DeleteOTP(email string) error {
	key := "otp:" + email
	return r.store.Delete(key)
}

func (r *Redis) SetOAuthState(state string, value []byte, exp time.Duration) error {
	key := "gstate:" + state
	return r.store.Set(key, value, exp)
}

func (r *Redis) GetOAuthState(state string) ([]byte, error) {
	key := "gstate:" + state
	val, err := r.store.Get(key)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *Redis) DeleteOAuthState(state string) error {
	key := "gstate:" + state
	return r.store.Delete(key)
}
