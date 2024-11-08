package entity

import (
	"errors"

	"github.com/lingwei0604/kitty/app/msg"
	"gorm.io/gorm"
)

var ErrRewardClaimed = errors.New(msg.RewardClaimed)
var ErrRelationCircled = errors.New("关系中不能有环")
var ErrRelationArgument = errors.New("错误的关系参数")
var ErrRelationSequence = errors.New("邀请者的注册日期晚于被邀请者")
var ErrRelationExist = errors.New("关系已经存在")
var ErrOrientationHasNotBeenCompleted = errors.New(msg.OrientationHasNotBeenCompleted)

type Relation struct {
	ID uint `gorm:"primaryKey;autoIncrementIncrement:2"`
	gorm.Model
	MasterID             uint `gorm:"index;"`
	ApprenticeID         uint `gorm:"index"`
	Master               User
	Apprentice           User
	Depth                int
	OrientationCompleted bool
	OrientationSteps     []OrientationStep
	RewardClaimed        bool
	PackageName          string
}

func NewRelation(apprentice *User, master *User, steps []OrientationStep) *Relation {
	return &Relation{
		MasterID:             master.ID,
		ApprenticeID:         apprentice.ID,
		Master:               *master,
		Apprentice:           *apprentice,
		Depth:                1,
		OrientationCompleted: len(steps) == 0,
		OrientationSteps:     steps,
		RewardClaimed:        false,
		PackageName:          apprentice.PackageName,
	}
}

func NewIndirectRelation(apprentice *User, master *User, steps []OrientationStep) *Relation {
	return &Relation{
		MasterID:             master.ID,
		ApprenticeID:         apprentice.ID,
		Master:               *master,
		Apprentice:           *apprentice,
		Depth:                2,
		OrientationCompleted: len(steps) == 0,
		OrientationSteps:     steps,
		RewardClaimed:        false,
		PackageName:          apprentice.PackageName,
	}
}

func (r *Relation) CompleteStep(step OrientationStep) {
	var orientationCompleted = true
	for n := range r.OrientationSteps {
		// update step status
		if r.OrientationSteps[n].Equals(step) {
			r.OrientationSteps[n].StepCompleted = true
		}
		// update orientationCompleted flag
		// 只要有一步没有完成，总的初始任务就没有完成
		if !r.OrientationSteps[n].StepCompleted {
			orientationCompleted = false
		}
	}
	r.OrientationCompleted = orientationCompleted
}

func (r *Relation) Validate() error {
	if r.MasterID == 0 {
		return ErrRelationArgument
	}
	if r.ApprenticeID == 0 {
		return ErrRelationArgument
	}
	if r.ApprenticeID == r.MasterID {
		return ErrRelationArgument
	}
	if r.Master.CreatedAt.After(r.Apprentice.CreatedAt) {
		return ErrRelationSequence
	}
	return nil
}

func (r *Relation) ClaimReward() error {
	if r.RewardClaimed {
		return ErrRewardClaimed
	}
	if !r.OrientationCompleted {
		return ErrOrientationHasNotBeenCompleted
	}
	r.RewardClaimed = true
	return nil
}

func (r *Relation) Connect(grandMaster *User, descendants []Relation) (addition []Relation, err error) {
	newRelations := []Relation{*r}
	if grandMaster != nil && grandMaster.ID != 0 {
		var steps = make([]OrientationStep, len(r.OrientationSteps))
		copy(steps, r.OrientationSteps)
		newRelations = append(newRelations, *NewIndirectRelation(&r.Apprentice, grandMaster, steps))
	}

	// 检测四阶环
	if circleDetected(&r.Master, grandMaster, descendants) {
		return nil, ErrRelationCircled
	}

	for _, descendant := range descendants {
		if descendant.Depth == 2 {
			continue
		}
		apprentice := User{Model: gorm.Model{ID: descendant.ApprenticeID}}
		var steps = make([]OrientationStep, len(r.OrientationSteps))
		copy(steps, r.OrientationSteps)
		newRelations = append(newRelations, *NewIndirectRelation(&apprentice, &r.Master, steps))
	}
	return newRelations, nil
}

func circleDetected(master, grandMaster *User, descendants []Relation) bool {
	if grandMaster != nil && grandMaster.ID != 0 {
		return in(grandMaster, descendants) || in(master, descendants)
	}
	return in(master, descendants)
}

func in(user *User, descendants []Relation) bool {
	for _, v := range descendants {
		if user.ID == v.ApprenticeID {
			return true
		}
	}
	return false
}

type OrientationStep struct {
	gorm.Model
	RelationID    uint `gorm:"index"`
	EventId       int
	ChineseName   string
	EventType     string
	StepCompleted bool
}

func (o OrientationStep) Equals(other OrientationStep) bool {
	return o.EventId == other.EventId && o.EventType == other.EventType
}
