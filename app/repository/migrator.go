package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lingwei0604/kitty/app/entity"
	"github.com/lingwei0604/kitty/pkg/contract"
	"gorm.io/gorm"
)

func ProvideMigrator(db *gorm.DB, appName contract.AppName) *gormigrate.Gormigrate {
	return gormigrate.New(db, &gormigrate.Options{
		TableName: fmt.Sprintf("%s_migrations", appName.String()),
	}, []*gormigrate.Migration{
		{
			ID: "202010280100",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					gorm.Model
					UserID    uint
					Os        uint8
					Imei      string
					Idfa      string
					Oaid      string
					Suuid     string
					Mac       string
					AndroidId string
					// 仅供数据库去重使用，应用不应依赖该字段，以免去重条件发生变化
					Hash string `gorm:"type:varchar(255);uniqueIndex:hash_index,sort:desc"`
				}
				type User struct {
					gorm.Model
					UserName      string         `json:"user_name" gorm:"default:test"`
					WechatOpenId  sql.NullString `json:"wechat_openid" gorm:"type:varchar(255);uniqueIndex:wechat_openid_index"`
					WechatUnionId sql.NullString `json:"wechat_unionid"`
					HeadImg       string         `json:"head_img"`
					Gender        int            `json:"gender"`
					Birthday      string         `json:"birthday" gorm:"default:2000-01-01"`
					Mobile        sql.NullString `json:"mobile" gorm:"type:varchar(255);uniqueIndex:mobile_index"`
					CommonSUUID   string         `json:"common_suuid"`
					Devices       []Device
					Channel       string `json:"channel"`
					VersionCode   string `json:"version_code"`
					InviteCode    string `json:"invite_code"`
					PackageName   string `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1"`
				}
				return db.AutoMigrate(
					&User{},
					&Device{},
				)
			},
			Rollback: func(db *gorm.DB) error {
				type User struct{}
				type Device struct{}
				return db.Migrator().DropTable(&User{}, &Device{})
			},
		},
		{
			ID: "202011030100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					ThirdPartyId string
				}
				if !db.Migrator().HasColumn(&User{}, "ThirdPartyId") {
					return db.Migrator().AddColumn(&User{}, "ThirdPartyId")
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					ThirdPartyId string
				}
				if db.Migrator().HasColumn(&entity.User{}, "ThirdPartyId") {
					return db.Migrator().DropColumn(&entity.User{}, "ThirdPartyId")
				}
				return nil
			},
		},
		{
			ID: "202011050100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					PackageName  string         `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1;uniqueIndex:taobao_openid_index,priority:1"`
					TaobaoOpenId sql.NullString `json:"taobao_openid" gorm:"type:varchar(255);uniqueIndex:taobao_openid_index"`
				}
				if !db.Migrator().HasColumn(&User{}, "TaobaoOpenId") {
					err := db.Migrator().AddColumn(&User{}, "TaobaoOpenId")
					if err != nil {
						return err
					}
					err = db.Migrator().CreateIndex(&User{}, "taobao_openid_index")
					if err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					PackageName  string         `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1;uniqueIndex:taobao_openid_index,priority:1"`
					TaobaoOpenId sql.NullString `json:"taobao_openid" gorm:"type:varchar(255);uniqueIndex:taobao_openid_index"`
				}
				if db.Migrator().HasColumn(&entity.User{}, "TaobaoOpenId") {
					err := db.Migrator().DropIndex(&User{}, "taobao_openid_index")
					if err != nil {
						return err
					}
					return db.Migrator().DropColumn(&User{}, "TaobaoOpenId")
				}
				return nil
			},
		},
		{
			ID: "202011130100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					WechatExtra []byte `gorm:"type:blob"`
					TaobaoExtra []byte `gorm:"type:blob"`
				}

				err := db.Migrator().AddColumn(&User{}, "WechatExtra")
				if err != nil {
					return err
				}
				err = db.Migrator().AddColumn(&User{}, "TaobaoExtra")
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					WechatExtra []byte `gorm:"type:blob"`
					TaobaoExtra []byte `gorm:"type:blob"`
				}
				err := db.Migrator().DropColumn(&User{}, "TaobaoExtra")
				if err != nil {
					return err
				}
				return db.Migrator().DropColumn(&User{}, "WechatExtra")

			},
		},
		{
			ID: "202011180100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					gorm.Model
				}
				type OrientationStep struct {
					gorm.Model
					RelationID    uint `gorm:"index"`
					Name          string
					StepCompleted bool
				}
				type Relation struct {
					gorm.Model
					MasterID             uint `gorm:"index"`
					ApprenticeID         uint `gorm:"index"`
					Master               User `gorm:"foreignKey:MasterID"`
					Apprentice           User `gorm:"foreignKey:ApprenticeID"`
					Depth                int
					OrientationCompleted bool
					OrientationSteps     []OrientationStep
					RewardClaimed        bool
				}

				return db.AutoMigrate(
					&OrientationStep{},
					&Relation{},
				)
			},
			Rollback: func(db *gorm.DB) error {
				type Relation struct{}
				type OrientationStep struct{}
				return db.Migrator().DropTable(&Relation{}, &OrientationStep{})
			},
		},
		{
			ID: "202012010100",
			Migrate: func(db *gorm.DB) error {
				type OrientationStep struct {
					gorm.Model
					RelationID    uint `gorm:"index"`
					EventId       int
					ChineseName   string
					EventType     string
					StepCompleted bool
				}
				db.Migrator().DropColumn(&OrientationStep{}, "Name")
				return db.AutoMigrate(
					&OrientationStep{},
				)
			},
			Rollback: func(db *gorm.DB) error {
				type Relation struct{}
				type OrientationStep struct{}
				db.Migrator().DropColumn(&OrientationStep{}, "EventId")
				db.Migrator().DropColumn(&OrientationStep{}, "ChineseName")
				db.Migrator().DropColumn(&OrientationStep{}, "EventType")
				db.Migrator().AddColumn(&OrientationStep{}, "Name")
				return nil
			},
		},
		{
			ID: "202012110100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					HeadImg string `gorm:"default:https://ad-static-xg.tagtic.cn/ad-material/file/0b8f18e1e666474291174ba316cccb51.png"`
				}
				return db.Migrator().AlterColumn(&User{}, "HeadImg")
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					HeadImg string
				}
				return db.Migrator().AlterColumn(&User{}, "HeadImg")
			},
		},
		{
			ID: "202012220100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					UserName string `json:"user_name" gorm:"default:test;type:varchar(30)"`
				}
				return db.Migrator().AlterColumn(&User{}, "UserName")
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					UserName string `json:"user_name" gorm:"default:test"`
				}
				return db.Migrator().AlterColumn(&User{}, "UserName")
			},
		},
		{
			ID: "202103240100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					PackageName string `gorm:"type:varchar(255);index:suuid_index,priority:1;uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1;uniqueIndex:taobao_openid_index,priority:1"`
					CommonSUUID string `gorm:"type:varchar(255);index:suuid_index,priority:2"`
				}
				if err := db.Migrator().AlterColumn(&User{}, "common_s_uuid"); err != nil {
					return err
				}
				return db.Migrator().CreateIndex(&User{}, "suuid_index")
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					PackageName string `gorm:"type:varchar(255);uniqueIndex:mobile_index,priority:1;uniqueIndex:wechat_openid_index,priority:1;uniqueIndex:taobao_openid_index,priority:1"`
					CommonSUUID string `gorm:""`
				}

				if err := db.Migrator().DropIndex(&User{}, "suuid_index"); err != nil {
					return err
				}
				if err := db.Migrator().AlterColumn(&User{}, "common_s_uuid"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202103260100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					CommonSMID sql.NullString `gorm:"type:varchar(255);"`
				}
				if err := db.Migrator().AddColumn(&User{}, "common_sm_id"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					CommonSMID sql.NullString
				}
				if err := db.Migrator().DropColumn(&User{}, "common_sm_id"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202103260200",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					SMID string `gorm:"type:varchar(255);"`
				}
				raw, _ := db.DB()
				raw.Exec("LOCK TABLES kitty_devices WRITE;")
				defer raw.Exec("UNLOCK TABLES;")
				if err := db.Migrator().AddColumn(&Device{}, "sm_id"); err != nil {
					return err
				}

				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type Device struct {
					SMID string `gorm:"type:varchar(255);"`
				}
				if err := db.Migrator().DropColumn(&Device{}, "sm_id"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202104070100",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					IP string `gorm:"type:varchar(255);"`
				}
				type User struct {
					CampaignID sql.NullString `gorm:"type:varchar(255);"`
					CID        sql.NullString `gorm:"type:varchar(255);"`
					AID        sql.NullString `gorm:"type:varchar(255);"`
				}
				raw, _ := db.DB()

				raw.Exec("LOCK TABLES kitty_devices WRITE;")
				defer raw.Exec("UNLOCK TABLES;")
				if err := db.Migrator().AddColumn(&Device{}, "ip"); err != nil {
					return err
				}
				raw.Exec("LOCK TABLES kitty_users WRITE;")
				if err := db.Migrator().AddColumn(&User{}, "campaign_id"); err != nil {
					return err
				}
				if err := db.Migrator().AddColumn(&User{}, "c_id"); err != nil {
					return err
				}
				if err := db.Migrator().AddColumn(&User{}, "a_id"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type Device struct {
					IP string
				}
				type User struct {
					CampaignID string
					CID        string
					AID        string
				}
				if err := db.Migrator().DropColumn(&Device{}, "ip"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&User{}, "campaign_id"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&User{}, "c_id"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&User{}, "a_id"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202106020100",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					Oaid string `gorm:"type:varchar(255);index:oaid"`
					Imei string `gorm:"type:varchar(255);index:imei"`
				}
				raw, _ := db.DB()
				raw.Exec("LOCK TABLES kitty_devices WRITE;")
				defer raw.Exec("UNLOCK TABLES;")
				if err := db.Migrator().AlterColumn(&Device{}, "oaid"); err != nil {
					return err
				}
				if err := db.Migrator().AlterColumn(&Device{}, "imei"); err != nil {
					return err
				}
				if err := db.Migrator().CreateIndex(&Device{}, "oaid"); err != nil {
					return err
				}
				if err := db.Migrator().CreateIndex(&Device{}, "imei"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type Device struct {
					Oaid string `gorm:"index:oaid"`
					Imei string `gorm:"index:imei"`
				}
				if err := db.Migrator().AlterColumn(&Device{}, "oaid"); err != nil {
					return err
				}
				if err := db.Migrator().AlterColumn(&Device{}, "imei"); err != nil {
					return err
				}
				if err := db.Migrator().DropIndex(&Device{}, "oaid"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&Device{}, "imei"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202107290100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					UnionSite string
				}
				if err := db.Migrator().AddColumn(&User{}, "union_site"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					UnionSite string
				}
				if err := db.Migrator().DropColumn(&User{}, "union_site"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202205131001",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					Idfa string `gorm:"type:varchar(255);index:idfa"`
				}
				raw, _ := db.DB()
				raw.Exec("LOCK TABLES kitty_devices WRITE;")
				defer raw.Exec("UNLOCK TABLES;")
				if err := db.Migrator().AlterColumn(&Device{}, "idfa"); err != nil {
					return err
				}
				if err := db.Migrator().CreateIndex(&Device{}, "idfa"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type Device struct {
					Idfa string `gorm:"type:varchar(255);index:idfa"`
				}
				if err := db.Migrator().DropColumn(&Device{}, "idfa"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202205161607",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					AndroidId string `gorm:"type:varchar(255);index:android_id"`
				}
				raw, _ := db.DB()
				raw.Exec("LOCK TABLES kitty_devices WRITE;")
				defer raw.Exec("UNLOCK TABLES;")
				if err := db.Migrator().AlterColumn(&Device{}, "android_id"); err != nil {
					return err
				}
				if err := db.Migrator().CreateIndex(&Device{}, "android_id"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type Device struct {
					AndroidId string `gorm:"type:varchar(255);index:android_id"`
				}
				if err := db.Migrator().DropColumn(&Device{}, "android_id"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202205161635",
			Migrate: func(db *gorm.DB) error {
				type Device struct {
					Suuid string `gorm:"type:varchar(255);index:suuid"`
				}
				raw, _ := db.DB()
				raw.Exec("LOCK TABLES kitty_devices WRITE;")
				defer raw.Exec("UNLOCK TABLES;")
				if err := db.Migrator().AlterColumn(&Device{}, "suuid"); err != nil {
					return err
				}
				if err := db.Migrator().CreateIndex(&Device{}, "suuid"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type Device struct {
					Suuid string `gorm:"type:varchar(255);index:suuid"`
				}
				if err := db.Migrator().DropColumn(&Device{}, "suuid"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "202206281825",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					CtaChannel sql.NullString `json:"cta_channel" gorm:"type:varchar(100)"`
				}
				if !db.Migrator().HasColumn(&User{}, "cta_channel") {
					err := db.Migrator().AddColumn(&User{}, "cta_channel")
					if err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type User struct {
					CtaChannel sql.NullString `json:"cta_channel" gorm:"type:varchar(100)"`
				}
				if db.Migrator().HasColumn(&entity.User{}, "cta_channel") {
					return db.Migrator().DropColumn(&User{}, "cta_channel")
				}
				return nil
			},
		},
		{
			ID: "202207050100",
			Migrate: func(db *gorm.DB) error {
				type Relation struct {
					PackageName string `json:"package_name" gorm:"type:varchar(255)"`
				}
				if !db.Migrator().HasColumn(&Relation{}, "package_name") {
					err := db.Migrator().AddColumn(&Relation{}, "package_name")
					if err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(db *gorm.DB) error {
				type Relation struct {
					PackageName string `json:"package_name" gorm:"type:varchar(255)"`
				}
				if db.Migrator().HasColumn(&entity.Relation{}, "package_name") {
					return db.Migrator().DropColumn(&Relation{}, "package_name")
				}
				return nil
			},
		},
	})
}
