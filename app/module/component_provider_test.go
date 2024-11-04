package module

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type IDAssignerImpl uint64

func (i *IDAssignerImpl) ID(ctx context.Context, packageName string, suuid string) (uint, error) {
	atomic.AddUint64((*uint64)(i), 1)
	return uint(atomic.LoadUint64((*uint64)(i))), nil
}

func TestProvideGormDB(t *testing.T) {
	sqlite := sqlite.Open(":memory:")
	db, _, _ := ProvideGormDB(sqlite, &gorm.Config{}, opentracing.GlobalTracer())

	for i := 0; i < 200; i++ {
		go func() {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Millisecond)
				db.Exec("SELECT 1")
			}
		}()
	}
	assert.Never(t, func() bool {
		sqlDB, err := db.DB()
		if err != nil {
			return true
		}
		t.Log("Open Connections: ", sqlDB.Stats().OpenConnections)
		return sqlDB.Stats().OpenConnections > maxOpenConns
	}, time.Second, time.Millisecond)
}
