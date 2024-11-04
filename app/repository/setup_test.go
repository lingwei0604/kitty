package repository

import (
	"context"
	"flag"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	mc "git.yingzhongshare.com/mkt/kitty/pkg/contract/mocks"
	"gorm.io/gorm/logger"

	"git.yingzhongshare.com/mkt/kitty/app/entity"
	"git.yingzhongshare.com/mkt/kitty/pkg/config"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var m *gormigrate.Gormigrate
var db *gorm.DB

var useMysql bool

func init() {
	flag.BoolVar(&useMysql, "mysql", false, "use local mysql for testing")
}

type mockID struct {
	i uint64
}

func (m *mockID) ID(ctx context.Context, packageName, suuid string) (uint, error) {
	atomic.AddUint64(&m.i, 1)
	return uint(atomic.LoadUint64(&m.i)), nil
}

func setUp(t *testing.T) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	var err error
	if !useMysql {
		db, err = gorm.Open(sqlite.Open(":memory:?cache=shared"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "kitty_", // 表名前缀，`User` 的表名应该是 `kitty_users`
			},
		})
	} else {
		db, err = gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/monetization?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "kitty_", // 表名前缀，`User` 的表名应该是 `kitty_users`
			},
		})
	}

	if err != nil {
		t.Fatal("failed to connect database")
	}
	db.Set("IDAssigner", &mockID{})
	m = ProvideMigrator(db, config.AppName("test"))
	err = m.Migrate()
	if err != nil {
		tearDown()
		t.Fatalf("failed migration: %s", err)
	}
}

func tearDown() {
	db.Migrator().DropTable(&entity.Device{}, &entity.Relation{}, &entity.User{}, &entity.OrientationStep{}, "test_migrations")
}

func user(id uint) entity.User {
	return entity.User{
		Model: gorm.Model{
			ID: id,
		},
	}
}

func getConf() contract.ConfigReader {
	conf := &mc.ConfigReader{}
	conf.On("String", "incrKey").Return("kitty-users-id", nil)
	return conf
}
