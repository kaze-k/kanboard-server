package global

import (
	"context"
	"fmt"
	"time"

	"server/internal/constant"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Ctx      context.Context
	client   *redis.Client
	duration time.Duration
}

var instance *RedisClient

var DURATION time.Duration = 30 * 24 * time.Hour

func getRedis() (*RedisClient, error) {
	redis := getInstance()
	redisClient := redis.client

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return redis, nil
}

func new() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     constant.RedisConfig.Addr,
		Password: constant.RedisConfig.Password,
		DB:       constant.RedisConfig.DB,
	})

	return &RedisClient{
		Ctx:      context.Background(),
		client:   client,
		duration: DURATION,
	}
}

func getInstance() *RedisClient {
	if instance == nil {
		instance = new()
	}
	return instance
}

func (r *RedisClient) Get(namespace string, key string) any {
	getKey := fmt.Sprintf("%s/%s", namespace, key)
	result, err := r.client.Get(context.Background(), getKey).Result()
	if err != nil {
		Logger.Error(err)
		return nil
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis Get", "key", getKey, "value", result)
	}

	return result
}

func (r *RedisClient) Set(namespace string, key string, value any, reset ...time.Duration) {
	expiration := r.duration
	if len(reset) > 0 {
		expiration = reset[0]
	}

	setKey := fmt.Sprintf("%s/%s", namespace, key)
	err := r.client.Set(context.Background(), setKey, value, expiration).Err()
	if err != nil {
		Logger.Error(err)
		panic(err)
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis Set", "key", setKey, "value", value)
	}
}

func (r *RedisClient) Delete(namespace string, key ...string) error {
	var deleteKeys []string
	for _, k := range key {
		deleteKeys = append(deleteKeys, fmt.Sprintf("%s/%s", namespace, k))
	}

	err := r.client.Del(context.Background(), deleteKeys...).Err()
	if err != nil {
		Logger.Error(err)
		return err
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis Delete", "keys", deleteKeys)
	}

	return err
}

func (r *RedisClient) GetTTL(namespace string, key string) time.Duration {
	ttlKey := fmt.Sprintf("%s/%s", namespace, key)
	ttl, err := r.client.TTL(context.Background(), ttlKey).Result()
	if err != nil {
		Logger.Error(err)
		return 0
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis GetTTL", "key", key, "ttl", ttl)
	}

	return ttl
}

func (r *RedisClient) HSet(namespace string, key string, value any) error {
	setKey := fmt.Sprintf("%s/%s", namespace, key)

	tx := r.client.TxPipeline()
	if err := r.client.HSet(r.Ctx, setKey, value).Err(); err != nil {
		Logger.Error(err)
	}
	if err := r.client.Expire(r.Ctx, setKey, r.duration).Err(); err != nil {
		Logger.Error(err)
	}
	if _, err := tx.Exec(r.Ctx); err != nil {
		Logger.Error(err)
		return err
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis HSet", "key", setKey, "value", value)
	}

	return nil
}

func (r *RedisClient) HGet(namespace string, key string, field string) string {
	getKey := fmt.Sprintf("%s/%s", namespace, key)

	result, getErr := r.client.HGet(r.Ctx, getKey, field).Result()
	if getErr != nil {
		Logger.Error(getErr)
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis HGet", "key", getKey, "value", result)
	}

	return result
}

func (r *RedisClient) HGetAll(namespace string, key string) map[string]string {
	getKey := fmt.Sprintf("%s/%s", namespace, key)
	result, err := r.client.HGetAll(r.Ctx, getKey).Result()
	if err != nil {
		Logger.Error(err)
		return nil
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis HGetAll", "key", getKey, "value", result)
	}

	return result
}

func (r *RedisClient) HDelete(namespace string, key string, fields ...string) error {
	deleteKey := fmt.Sprintf("%s/%s", namespace, key)

	err := r.client.HDel(r.Ctx, deleteKey, fields...).Err()
	if err != nil {
		Logger.Error(err)
		return err
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis HDelete", "keys", deleteKey)
	}

	return nil
}

func (r *RedisClient) SAdd(namespace string, key string, value any, expire bool) error {
	setKey := fmt.Sprintf("%s/%s", namespace, key)

	tx := r.client.TxPipeline()
	if err := r.client.SAdd(r.Ctx, setKey, value).Err(); err != nil {
		Logger.Error(err)
	}
	if expire {
		if err := r.client.Expire(r.Ctx, setKey, r.duration).Err(); err != nil {
			Logger.Error(err)
		}
	}
	_, err := tx.Exec(r.Ctx)
	if err != nil {
		Logger.Error(err)
		return err
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis SAdd", "key", setKey, "value", value)
	}

	return nil
}

func (r *RedisClient) SRem(namespace string, key string, value any) error {
	setKey := fmt.Sprintf("%s/%s", namespace, key)
	err := r.client.SRem(r.Ctx, setKey, value).Err()
	if err != nil {
		Logger.Error(err)
		return err
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis SRem", "key", setKey, "value", value)
	}

	return err
}

func (r *RedisClient) SMembers(namespace string, key string) []string {
	setKey := fmt.Sprintf("%s/%s", namespace, key)
	result, err := r.client.SMembers(r.Ctx, setKey).Result()
	if err != nil {
		Logger.Error(err)
		return nil
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis SMembers", "key", setKey, "value", result)
	}

	return result
}

func (r *RedisClient) SCard(namespace string, key string) int64 {
	setKey := fmt.Sprintf("%s/%s", namespace, key)
	result, err := r.client.SCard(r.Ctx, setKey).Result()
	if err != nil {
		Logger.Error(err)
		return 0
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis SCard", "key", setKey, "value", result)
	}

	return result
}

func (r *RedisClient) Publish(channel string, content string) error {
	err := r.client.Publish(r.Ctx, channel, content).Err()
	if err != nil {
		Logger.Error(err)
		return err
	}

	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis Publish", "channel", channel, "content", content)
	}

	return nil
}

func (r *RedisClient) Subscribe(channel string) *redis.PubSub {
	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis Subscribe", "channel", channel)
	}

	return r.client.Subscribe(r.Ctx, channel)
}

func (r *RedisClient) Scan(namespace string, match string, count int64) *redis.ScanIterator {
	getMatch := fmt.Sprintf("%s/%s", namespace, match)
	if constant.EnvConfig.Mode == "debug" {
		Logger.Infow("redis Scan", "match", getMatch, "count", count)
	}
	return r.client.Scan(r.Ctx, 0, getMatch, count).Iterator()
}
