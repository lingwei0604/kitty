package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/lingwei0604/kitty/app/entity"
	code "github.com/lingwei0604/kitty/pkg/invitecode"
	"github.com/lingwei0604/kitty/pkg/kjwt"
	"github.com/lingwei0604/kitty/share/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getConf() *ShareConfig {
	return &ShareConfig{
		OrientationEvents: []OrientationEvent{{
			Id:   0,
			Type: "task",
		}},
		Url: "http://www.donews.com?%s",
		Reward: struct {
			Level1 int `yaml:"level1"`
			Level2 int `yaml:"level2"`
		}{
			100,
			50,
		},
		TaskId: "111",
	}
}

type MockClient func(ctx context.Context, dto *XTaskRequest) (*XTaskResponse, error)

func (m MockClient) Request(ctx context.Context, dto *XTaskRequest) (*XTaskResponse, error) {
	return m(ctx, dto)
}

func TestInvitationManager_AddToken(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		service  InvitationManager
		userId   uint
		masterId uint
		out      error
	}{
		{
			"插入token",
			InvitationManager{
				conf: getConf(),
				rr: (func() RelationRepository {
					ur := &mocks.RelationRepository{}
					ur.On("AddRelations", mock.Anything, mock.Anything).Return(nil).Once()
					return ur
				})(),
				tokenizer: code.NewTokenizer("foo"),
			},

			1,
			2,
			nil,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			apprentice := user(cc.userId)
			master := user(cc.masterId)
			err := cc.service.AddToken(context.Background(), &apprentice, &master)
			assert.True(t, errors.Is(err, cc.out))
		})
	}
}

func TestInvitationManager_ClaimReward(t *testing.T) {
	t.Parallel()
	rels := []entity.Relation{{
		MasterID:             1,
		ApprenticeID:         2,
		Master:               user(1),
		Apprentice:           user(2),
		Depth:                1,
		OrientationCompleted: true,
		OrientationSteps:     nil,
		RewardClaimed:        false,
	}}
	cases := []struct {
		name         string
		service      InvitationManager
		masterId     uint64
		apprenticeId uint64
		out          error
	}{
		{
			"正常claim",
			InvitationManager{
				conf: getConf(),
				rr: (func() RelationRepository {
					ur := &mocks.RelationRepository{}
					ur.On("UpdateRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, apprentice *entity.User, existingRelationCallback func(relations []entity.Relation) error) error {
						return existingRelationCallback(rels)
					}).Once()
					return ur
				})(),
				tokenizer: code.NewTokenizer("foo"),
				xtaskClient: MockClient(func(ctx context.Context, dto *XTaskRequest) (*XTaskResponse, error) {
					return &XTaskResponse{Code: 0}, nil
				}),
			},

			1,
			2,
			nil,
		},
		{
			"因为HTTP请求失败所以无法完成",
			InvitationManager{
				conf: getConf(),
				rr: (func() RelationRepository {
					ur := &mocks.RelationRepository{}
					ur.On("UpdateRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, apprentice *entity.User, existingRelationCallback func(relations []entity.Relation) error) error {
						return existingRelationCallback([]entity.Relation{{
							MasterID:             1,
							ApprenticeID:         2,
							Master:               user(1),
							Apprentice:           user(2),
							Depth:                1,
							OrientationCompleted: true,
							OrientationSteps:     nil,
							RewardClaimed:        false,
						}})
					}).Once()
					return ur
				})(),
				tokenizer: code.NewTokenizer("foo"),
				xtaskClient: MockClient(func(ctx context.Context, dto *XTaskRequest) (*XTaskResponse, error) {
					return &XTaskResponse{Code: 2, Msg: "foo"}, ErrFailedXtaskRequest
				}),
			},

			1,
			2,
			ErrFailedXtaskRequest,
		},
		{
			"由于OrientationStep未完成所以无法claim",
			InvitationManager{
				conf: getConf(),
				rr: (func() RelationRepository {
					ur := &mocks.RelationRepository{}
					ur.On("UpdateRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, apprentice *entity.User, existingRelationCallback func(relations []entity.Relation) error) error {
						return existingRelationCallback([]entity.Relation{{
							MasterID:             1,
							ApprenticeID:         2,
							Master:               user(1),
							Apprentice:           user(2),
							Depth:                1,
							OrientationCompleted: false,
							OrientationSteps:     nil,
							RewardClaimed:        false,
						}})
					}).Once()
					return ur
				})(),
				tokenizer: code.NewTokenizer("foo"),
			},

			1,
			2,
			entity.ErrOrientationHasNotBeenCompleted,
		},
		{
			"由于关系不存在所以无法claim",
			InvitationManager{
				conf: getConf(),
				rr: (func() RelationRepository {
					ur := &mocks.RelationRepository{}
					ur.On("UpdateRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, apprentice *entity.User, existingRelationCallback func(relations []entity.Relation) error) error {
						return existingRelationCallback([]entity.Relation{})
					}).Once()
					return ur
				})(),
				tokenizer: code.NewTokenizer("foo"),
			},

			1,
			2,
			ErrNoRewardAvailable,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			err := cc.service.ClaimReward(context.Background(), cc.masterId, cc.apprenticeId)
			assert.True(t, errors.Is(err, cc.out))
			if err == nil {
				assert.True(t, rels[0].RewardClaimed)
			}

		})
	}
}

func TestInvitationManager_AdvanceStep(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name         string
		service      InvitationManager
		apprenticeId uint64
		eventName    ReceivedEvent
		out          error
	}{
		{
			"正常前进",
			InvitationManager{
				logger: log.NewNopLogger(),
				conf:   getConf(),
				rr: (func() RelationRepository {
					ur := &mocks.RelationRepository{}
					ur.On("UpdateRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, apprentice *entity.User, existingRelationCallback func(relations []entity.Relation) error) error {
						return existingRelationCallback([]entity.Relation{{
							MasterID:             1,
							ApprenticeID:         2,
							Master:               user(1),
							Apprentice:           user(2),
							Depth:                1,
							OrientationCompleted: false,
							OrientationSteps: []entity.OrientationStep{{
								EventId:       1,
								StepCompleted: false,
							}},
							RewardClaimed: false,
						}})
					}).Once()
					return ur
				})(),
				tokenizer: code.NewTokenizer("foo"),
			},

			1,
			ReceivedEvent{
				Id:   0,
				Type: "task",
			},
			nil,
		},
		{
			"关系不存在时静默处理",
			InvitationManager{
				conf: getConf(),
				rr: (func() RelationRepository {
					ur := &mocks.RelationRepository{}
					ur.On("UpdateRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, apprentice *entity.User, existingRelationCallback func(relations []entity.Relation) error) error {
						return existingRelationCallback([]entity.Relation{})
					}).Once()
					return ur
				})(),
				tokenizer: code.NewTokenizer("foo"),
				logger:    log.NewNopLogger(),
			},

			1,
			ReceivedEvent{
				Id:   0,
				Type: "task",
			},
			nil,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			err := cc.service.CompleteStep(context.Background(), cc.apprenticeId, cc.eventName)
			assert.True(t, errors.Is(err, cc.out))
		})
	}
}

func TestInvitationManager_GetUrl(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		service InvitationManager
		ctx     context.Context
		out     string
	}{
		{
			"拼接URL",
			InvitationManager{
				conf:      getConf(),
				rr:        nil,
				tokenizer: code.NewTokenizer("foo"),
				logger:    log.NewNopLogger(),
			},
			context.WithValue(context.Background(), jwt.JWTClaimsContextKey, &kjwt.Claim{
				PackageName: "com.donews.www",
				UserId:      100,
				Channel:     "foo",
				VersionCode: "1000",
			}),
			"http://www.donews.com?channel=foo_share&invite_code=87V6lEZJvN&package_name=com.donews.www&user_id=100&version_code=1000",
		},
		{
			"二次分享拼接URL",
			InvitationManager{
				conf:      getConf(),
				rr:        nil,
				tokenizer: code.NewTokenizer("foo"),
				logger:    log.NewNopLogger(),
			},
			context.WithValue(context.Background(), jwt.JWTClaimsContextKey, &kjwt.Claim{
				PackageName: "com.donews.www",
				UserId:      100,
				Channel:     "foo_share",
				VersionCode: "1000",
			}),
			"http://www.donews.com?channel=foo_share&invite_code=87V6lEZJvN&package_name=com.donews.www&user_id=100&version_code=1000",
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			uri := cc.service.GetUrl(cc.ctx, kjwt.ClaimFromContext(cc.ctx))
			assert.Equal(t, cc.out, uri)
		})
	}
}

func TestInvitationManager_ListApprentice(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		service  InvitationManager
		masterId uint
		depth    int
		amount   int
	}{
		{
			"一级邀请列表",
			InvitationManager{
				conf: getConf(),
				rr: func() RelationRepository {
					var ur mocks.RelationRepository
					ur.On("QueryRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, condition entity.Relation) []entity.Relation {
						return []entity.Relation{{
							MasterID:             condition.MasterID,
							ApprenticeID:         1,
							Depth:                condition.Depth,
							OrientationCompleted: false,
							OrientationSteps:     nil,
							RewardClaimed:        false,
						}}
					}, nil).Once()
					return &ur
				}(),
				tokenizer: code.NewTokenizer("foo"),
			},
			1,
			1,
			100,
		},
		{
			"二级邀请列表",
			InvitationManager{
				conf: getConf(),
				rr: func() RelationRepository {
					var ur mocks.RelationRepository
					ur.On("QueryRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, condition entity.Relation) []entity.Relation {
						return []entity.Relation{{
							MasterID:             condition.MasterID,
							ApprenticeID:         1,
							Depth:                condition.Depth,
							OrientationCompleted: false,
							OrientationSteps:     nil,
							RewardClaimed:        false,
						}}
					}, nil).Once()
					return &ur
				}(),
				tokenizer: code.NewTokenizer("foo"),
			},
			10,
			2,
			50,
		},
		{
			"多个返回值",
			InvitationManager{
				conf: getConf(),
				rr: func() RelationRepository {
					var ur mocks.RelationRepository
					ur.On("QueryRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, condition entity.Relation) []entity.Relation {
						return []entity.Relation{{
							MasterID:             condition.MasterID,
							ApprenticeID:         1,
							Depth:                condition.Depth,
							OrientationCompleted: false,
							OrientationSteps:     nil,
							RewardClaimed:        false,
						}, {
							MasterID:             condition.MasterID,
							ApprenticeID:         2,
							Depth:                condition.Depth,
							OrientationCompleted: false,
							OrientationSteps:     nil,
							RewardClaimed:        false,
						}}
					}, nil).Once()
					return &ur
				}(),
				tokenizer: code.NewTokenizer("foo"),
			},
			1,
			1,
			100,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			rel, err := cc.service.ListApprentices(context.Background(), uint64(cc.masterId), cc.depth)
			assert.NoError(t, err)
			for i, r := range rel {
				assert.Equal(t, cc.masterId, r.MasterID)
				assert.Equal(t, uint(1+i), r.ApprenticeID)
				assert.Equal(t, cc.depth, r.Depth)
				assert.Equal(t, cc.amount, r.Amount)
			}
		})
	}
}

func TestInvitationManager_ListMaster(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		service         InvitationManager
		apprenticeID    uint
		expectedMasters int
	}{
		{
			"两个师傅都在",
			InvitationManager{
				conf: getConf(),
				rr: func() RelationRepository {
					var ur mocks.RelationRepository
					ur.On("QueryRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, condition entity.Relation) []entity.Relation {
						return []entity.Relation{{
							MasterID:     1,
							Master:       user(1),
							ApprenticeID: condition.ApprenticeID,
							Depth:        1,
						}, {
							MasterID:     2,
							Master:       user(2),
							ApprenticeID: condition.ApprenticeID,
							Depth:        2,
						}}
					}, nil).Once()
					return &ur
				}(),
				tokenizer: code.NewTokenizer("foo"),
			},
			3,
			2,
		},
		{
			"只有一个师傅了",
			InvitationManager{
				conf: getConf(),
				rr: func() RelationRepository {
					var ur mocks.RelationRepository
					ur.On("QueryRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, condition entity.Relation) []entity.Relation {
						return []entity.Relation{{
							MasterID:     1,
							Master:       user(1),
							ApprenticeID: condition.ApprenticeID,
							Depth:        1,
						}}
					}, nil).Once()
					return &ur
				}(),
				tokenizer: code.NewTokenizer("foo"),
			},
			3,
			1,
		},
		{
			"没有师傅",
			InvitationManager{
				conf: getConf(),
				rr: func() RelationRepository {
					var ur mocks.RelationRepository
					ur.On("QueryRelations", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, condition entity.Relation) []entity.Relation {
						return []entity.Relation{}
					}, nil).Once()
					return &ur
				}(),
				tokenizer: code.NewTokenizer("foo"),
			},
			3,
			0,
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			master, grandMaster, err := cc.service.ListMaster(context.Background(), uint64(cc.apprenticeID))
			assert.NoError(t, err)
			if cc.expectedMasters >= 2 {
				assert.NotNil(t, master)
				assert.NotNil(t, grandMaster)
			}
			if cc.expectedMasters == 1 {
				assert.NotNil(t, master)
				assert.Nil(t, grandMaster)
			}
			if cc.expectedMasters == 0 {
				assert.Nil(t, master)
				assert.Nil(t, grandMaster)
			}
		})
	}
}
