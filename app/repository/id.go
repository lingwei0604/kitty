package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/lingwei0604/kitty/app/entity"
	"github.com/lingwei0604/kitty/pkg/contract"
)

const insertScript = `
local suuid = ARGV[1]
if suuid == ""
then
	return redis.call("incr", KEYS[1])
end
local id = redis.call("get", KEYS[2]..":"..suuid)
if id ~= false
then
	redis.call("expire", KEYS[2]..":"..suuid, "604800")
	return id
end
local result = redis.call("incr", KEYS[1])
redis.call("set", KEYS[2]..":"..suuid, result, "EX", "604800")
return result
`

const lookupScript = `
local suuid = ARGV[1]
if suuid == ""
then
	return redis.call("incr", KEYS[1])
end
local id = redis.call("get", KEYS[2]..":"..suuid)
if id ~= false
then
	return id
end
return redis.call("incr", KEYS[1])
`

type UniqueID struct {
	redis redis.UniversalClient
	key   string
}

func NewUniqueID(redis redis.UniversalClient, conf contract.ConfigReader) *UniqueID {
	return &UniqueID{
		redis: redis,
		key:   conf.String("incrKey"),
	}
}

func (u *UniqueID) ID(ctx context.Context, packageName string, suuid string) (uint, error) {
	id, err := u.redis.Eval(ctx, lookupScript, []string{u.key, fmt.Sprintf(
		"{%s}:preallocate:suuid",
		u.key,
	)}, strings.Join([]string{packageName, suuid}, ":")).Result()
	if err != nil {
		return 0, errors.Wrap(err, "cannot create id from redis")
	}
	return cast.ToUint(id), nil
}

func (u *UniqueID) Preallocate(ctx context.Context, packageName string, suuid string) (uint, error) {
	data, err := u.redis.Eval(ctx, insertScript, []string{u.key, fmt.Sprintf(
		"{%s}:preallocate:suuid",
		u.key,
	)}, strings.Join([]string{packageName, suuid}, ":")).Result()

	if err != nil {
		return 0, fmt.Errorf("preallocate id in redis: %w", err)
	}
	return cast.ToUint(data), nil
}

func (u *UniqueID) ClearPreallocate(ctx context.Context, packageName string, suuid string) error {
	return u.redis.Del(ctx, fmt.Sprintf(
		"{%s}:preallocate:suuid:%s",
		u.key, strings.Join([]string{packageName, suuid}, ":"),
	)).Err()
}

func (u *UniqueID) SetIDByDevice(ctx context.Context, packageName string, device *entity.Device, ID uint) error {
	pipe := u.redis.Pipeline()
	if device.Suuid != "" {
		pipe.Set(ctx, u.genKey("suuid", packageName, device.Suuid), ID, 30*86400*time.Second)
	}
	if device.GetOaid() != "" {
		pipe.Set(ctx, u.genKey("oaid", packageName, device.GetOaid()), ID, 30*86400*time.Second)
	}
	if device.GetAndroidId() != "" {
		pipe.Set(ctx, u.genKey("android", packageName, device.GetAndroidId()), ID, 30*86400*time.Second)
	}
	if device.GetIdfa() != "" {
		pipe.Set(ctx, u.genKey("idfa", packageName, device.GetIdfa()), ID, 30*86400*time.Second)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (u *UniqueID) GetIDByDevice(ctx context.Context, packageName string, device *entity.Device) (uint, error) {
	var keys []string
	if device.Suuid != "" {
		keys = append(keys, u.genKey("suuid", packageName, device.Suuid))
	}
	if device.GetOaid() != "" {
		keys = append(keys, u.genKey("oaid", packageName, device.GetOaid()))
	}
	if device.GetAndroidId() != "" {
		keys = append(keys, u.genKey("android", packageName, device.GetAndroidId()))
	}
	if device.GetIdfa() != "" {
		keys = append(keys, u.genKey("idfa", packageName, device.GetIdfa()))
	}
	if len(keys) == 0 {
		return 0, ErrRecordNotFound
	}

	for _, key := range keys {
		res, err := u.redis.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return 0, err
		}
		if ID, err := strconv.ParseUint(res, 10, 64); err == nil && ID > 0 {
			return uint(ID), nil
		}
	}
	return 0, ErrRecordNotFound
}
func (u *UniqueID) DeleteCacheByDevice(ctx context.Context, packageName string, device *entity.Device) {
	var keys []string
	if device.Suuid != "" {
		keys = append(keys, u.genKey("suuid", packageName, device.Suuid))
	}
	if device.GetOaid() != "" {
		keys = append(keys, u.genKey("oaid", packageName, device.GetOaid()))
	}
	if device.GetAndroidId() != "" {
		keys = append(keys, u.genKey("android", packageName, device.GetAndroidId()))
	}
	if device.GetIdfa() != "" {
		keys = append(keys, u.genKey("idfa", packageName, device.GetIdfa()))
	}
	if len(keys) == 0 {
		return
	}
	pipe := u.redis.Pipeline()
	defer pipe.Close()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}
	pipe.Exec(ctx)
	return
}

func (u *UniqueID) SetRegisterTimeById(ctx context.Context, ID uint, registerTime int64) error {
	return u.redis.Set(ctx, strings.Join([]string{u.key, "register", "time", strconv.Itoa(int(ID))}, ":"), registerTime, 30*24*time.Hour).Err()
}

func (u *UniqueID) GetRegisterTimeById(ctx context.Context, ID uint) (int64, error) {
	res, err := u.redis.Get(ctx, strings.Join([]string{u.key, "register", "time", strconv.Itoa(int(ID))}, ":")).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	reg, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		return 0, err
	}
	return reg, nil
}

func (u *UniqueID) genKey(field, packageName, suffix string) string {
	return strings.Join([]string{u.key, "preallocate:device", packageName, suffix}, ":")
}
