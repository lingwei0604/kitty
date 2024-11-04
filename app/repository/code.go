package repository

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lingwei0604/kitty/app/msg"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/otredis"
	"github.com/pkg/errors"
)

const CodeKey = "CodeRepo"
const defaultTtl = 5 * time.Minute
const defaultRate = time.Minute

var ErrTooFrequent = errors.New(msg.ErrorTooFrequent)

type CodeRepo struct {
	client redis.Cmdable
	km     contract.Keyer
	ttl    time.Duration
	rate   time.Duration
	env    contract.Env
}

func NewCodeRepo(cmdable redis.Cmdable, keyer contract.Keyer, env contract.Env) *CodeRepo {
	return &CodeRepo{cmdable, otredis.With(keyer, CodeKey), defaultTtl, defaultRate, env}
}

func (c *CodeRepo) AddCode(ctx context.Context, mobile string) (code string, err error) {
	// 限制每个号码每分钟最多重新生成一个
	left, err := c.client.TTL(ctx, c.km.Key(mobile)).Result()
	if err != nil && err != redis.Nil {
		return "", errors.Wrap(err, "cannot connect to redis")
	}
	if left > c.ttl-c.rate {
		return "", ErrTooFrequent
	}
	n := rand.Intn(1_000_000)
	code = pad(n)
	_, err = c.client.Set(ctx, c.km.Key(mobile), code, c.ttl).Result()
	if err != nil {
		return "", errors.Wrap(err, "cannot persist code in redis")
	}
	return code, nil
}

func (c *CodeRepo) CheckCode(ctx context.Context, mobile, code string) (bool, error) {
	// 上架时提供给审核方登录使用
	if mobile == "18633697850" && code == "666666" {
		return true, nil
	}
	if !c.env.IsProd() && code == "666666" {
		return true, nil
	}
	value, err := c.client.Get(ctx, c.km.Key(mobile)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "cannot query code in redis")
	}
	return value == code, nil
}

func (c *CodeRepo) DeleteCode(ctx context.Context, mobile string) (err error) {
	_, err = c.client.Del(ctx, c.km.Key(mobile)).Result()
	return err
}

func pad(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}
