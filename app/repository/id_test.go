package repository

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"git.yingzhongshare.com/mkt/kitty/app/entity"
)

func TestPreallocate(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Skip("Set REDIS_ADDR to run Preallocate test")
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: strings.Split(os.Getenv("REDIS_ADDR"), ",")})
	idRepo := &UniqueID{
		redis: client,
		key:   "test",
	}
	id, err := idRepo.Preallocate(context.Background(), "", "foo")
	assert.NoError(t, err)
	t.Log(id)
}

func TestID(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Skip("Set REDIS_ADDR to run Preallocate test")
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: strings.Split(os.Getenv("REDIS_ADDR"), ",")})
	idRepo := &UniqueID{
		redis: client,
		key:   "test",
	}

	t.Run("without preallocate", func(t *testing.T) {
		id1, err := idRepo.ID(context.Background(), "aaa", "foo")
		assert.NoError(t, err)
		t.Log(id1)
		id2, _ := idRepo.ID(context.Background(), "aaa", "foo")
		assert.NoError(t, err)
		t.Log(id2)
		assert.NotEqual(t, id1, id2)
	})

	t.Run("with preallocate", func(t *testing.T) {
		idRepo.Preallocate(context.Background(), "", "bar")
		id1, _ := idRepo.ID(context.Background(), "", "bar")
		id2, _ := idRepo.ID(context.Background(), "", "bar")
		assert.Equal(t, id1, id2)
	})

}

func TestCacheIDWithDevice(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" {
		t.Skip("Set REDIS_ADDR to run Preallocate test")
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: strings.Split(os.Getenv("REDIS_ADDR"), ",")})
	idRepo := &UniqueID{
		redis: client,
		key:   "test",
	}
	ctx := context.Background()
	err := idRepo.SetIDByDevice(ctx, "test", &entity.Device{
		AndroidId: "1",
	}, 1)
	assert.NoError(t, err)

	id, err := idRepo.GetIDByDevice(ctx, "test", &entity.Device{
		AndroidId: "1",
	})
	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)

	err = idRepo.SetRegisterTimeById(ctx, 1, 1)
	assert.NoError(t, err)
	reg, err := idRepo.GetRegisterTimeById(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), reg)
}
