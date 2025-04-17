package captcha

import (
	"strings"
	"time"

	"server/internal/global"

	"github.com/mojocn/base64Captcha"
)

var source string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RedisStore struct {
	Namespace   string
	Duration    time.Duration
	RedisClient *global.RedisClient
}

func NewRedisStore(namespace string, duration time.Duration, redisClient *global.RedisClient) *RedisStore {
	return &RedisStore{
		Namespace:   namespace,
		Duration:    duration,
		RedisClient: redisClient,
	}
}

func (r *RedisStore) Get(id string, clear bool) string {
	result := r.RedisClient.Get(r.Namespace, id)
	if result != nil {
		return result.(string)
	}
	if clear {
		r.RedisClient.Delete(r.Namespace, id)
	}
	return ""
}

func (r *RedisStore) Set(id, value string) error {
	r.RedisClient.Set(r.Namespace, id, value, r.Duration)
	return nil
}

func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	result := r.Get(id, clear)
	return strings.ToLower(result) == strings.ToLower(answer)
}

type Captcha struct {
	Store *RedisStore
}

func NewCaptcha(store *RedisStore) *Captcha {
	return &Captcha{
		Store: store,
	}
}

func (c *Captcha) Generate() (string, string, string, error) {
	driver := base64Captcha.NewDriverString(
		80,
		240,
		2,
		base64Captcha.OptionShowHollowLine,
		4,
		source,
		nil,
		nil,
		nil,
	)
	captcha := base64Captcha.NewCaptcha(driver, c.Store)
	return captcha.Generate()
}

func (c *Captcha) Verify(id, answer string) bool {
	return c.Store.Verify(id, answer, true)
}
