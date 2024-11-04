package repository

import (
	"context"
	"encoding/json"
	"strconv"

	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/otredis"
	"github.com/go-redis/redis/v8"

	"time"
)

type UserCache struct {
	redis redis.UniversalClient
	key   contract.Keyer
}

func NewUserCache(redis redis.UniversalClient) *UserCache {
	return &UserCache{
		redis: redis,
		key:   otredis.NewKeyManager(":", "kitty"),
	}
}

func (u *UserCache) CacheID(ctx context.Context, id uint) error {
	return u.redis.Set(ctx, u.key.Key("uid", strconv.Itoa(int(id))), 1, 30*86400*time.Second).Err()
}

type BindAdCallback struct {
	Id         uint64 `json:"id,omitempty"`
	CampaignId string `json:"campaign_id,omitempty"`
	Cid        string `json:"cid,omitempty"`
	Aid        string `json:"aid,omitempty"`
	UnionSite  string `json:"union_site,omitempty"`
	CtaChannel string `json:"cta_channel,omitempty"`
}

func (u *UserCache) CacheBindAd(ctx context.Context, cb *BindAdCallback) error {
	bs, _ := json.Marshal(cb)
	return u.redis.Set(ctx, u.key.Key("bindad", strconv.Itoa(int(cb.Id))), string(bs), 3*24*time.Hour).Err()
}

func (u *UserCache) GetBindAd(ctx context.Context, id uint) (*BindAdCallback, error) {
	r, err := u.redis.Get(ctx, u.key.Key("bindad", strconv.Itoa(int(id)))).Result()
	if err != nil {
		return nil, err
	}
	var cb *BindAdCallback
	err = json.Unmarshal([]byte(r), &cb)
	if err != nil {
		return nil, err
	}

	return cb, nil
}
