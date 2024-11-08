package repository

import (
	"github.com/lingwei0604/kitty/app/entity"
	"github.com/lingwei0604/kitty/pkg/otgorm"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

func createUser(db *gorm.DB, name string) error {
	return db.Create(&entity.User{
		UserName:    name,
		CommonSUUID: xid.New().String(),
	}).Error
}

func createRelation(db *gorm.DB, masterId int, apprenticeId int) error {
	return db.Create(&entity.Relation{
		MasterID:             uint(masterId),
		ApprenticeID:         uint(apprenticeId),
		Depth:                1,
		OrientationCompleted: false,
		OrientationSteps: []entity.OrientationStep{
			{
				EventType:   "task",
				EventId:     1,
				ChineseName: "一个任务的中文名字",
			},
		},
		RewardClaimed: false,
	}).Error
}

func ProvideSeeder(db *gorm.DB) otgorm.Seeds {
	return otgorm.Seeds{
		Db: db,
		Seeds: []otgorm.Seed{
			{
				Name: "create 100 users",
				Run: func(db *gorm.DB) error {
					for i := 0; i < 100; i++ {
						err := createUser(db, "mr. Fake")
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name: "make all odd id masters; all even id apprentice",
				Run: func(db *gorm.DB) error {
					for i := 1; i < 100; i += 2 {
						err := createRelation(db, i, i+1)
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
		},
	}
}
