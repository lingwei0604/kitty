package repository

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestUserCache_CacheID(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Skip("set REDIS_ADDR to run UserCache test")
	}
	cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: strings.Split(os.Getenv("REDIS_ADDR"), ","),
	})

	ctx := context.TODO()

	if _, err := cli.Ping(ctx).Result(); err != nil {
		t.Fatal(err)
	}
	c := NewUserCache(cli)

	if err := c.CacheID(ctx, 1); err != nil {
		t.Fatal(err)
	}

	v, err := cli.Exists(ctx, "kitty:uid:1").Result()
	if err != nil {
		t.Fatal(err)
	}
	if v == 0 {
		t.Error("want key exists, but not exists")
	}

}

func TestUserCache_CacheBindAd(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Skip("set REDIS_ADDR to run UserCache test")
	}
	cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: strings.Split(os.Getenv("REDIS_ADDR"), ","),
	})

	ctx := context.TODO()

	if _, err := cli.Ping(ctx).Result(); err != nil {
		t.Fatal(err)
	}
	c := NewUserCache(cli)

	if err := c.CacheBindAd(ctx, &BindAdCallback{
		Id:         1,
		CampaignId: "1",
		Cid:        "1",
		Aid:        "1",
		UnionSite:  "1",
	}); err != nil {
		t.Fatal(err)
	}

	v, err := c.GetBindAd(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	if v.Aid != "1" {
		t.Errorf("want aid is 1, but get %s", v.Aid)
	}

}
