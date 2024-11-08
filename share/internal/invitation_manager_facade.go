package internal

import (
	"context"
	"github.com/lingwei0604/kitty/app/entity"
	jwt2 "github.com/lingwei0604/kitty/pkg/kjwt"

	"github.com/lingwei0604/kitty/pkg/config"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/pkg/errors"
)

type InvitationManagerFacade struct {
	Name    contract.AppName
	Factory InvitationManagerFactory
	DynConf config.DynamicConfigReader
}

func (im *InvitationManagerFacade) ListMaster(ctx context.Context, apprenticeId uint64) (master *entity.User, grandMaster *entity.User, err error) {
	m, err := im.getManager(ctx)
	if err != nil {
		return nil, nil, err
	}
	return m.ListMaster(ctx, apprenticeId)
}

func (im *InvitationManagerFacade) AddToken(ctx context.Context, apprentice, master *entity.User) error {
	m, err := im.getManager(ctx)
	if err != nil {
		return err
	}
	return m.AddToken(ctx, apprentice, master)
}

func (im *InvitationManagerFacade) ClaimReward(ctx context.Context, masterId uint64, apprenticeId uint64) error {
	m, err := im.getManager(ctx)
	if err != nil {
		return err
	}
	return m.ClaimReward(ctx, masterId, apprenticeId)
}

func (im *InvitationManagerFacade) CompleteStep(ctx context.Context, apprenticeId uint64, event ReceivedEvent) error {
	m, err := im.getManager(ctx)
	if err != nil {
		return err
	}
	return m.CompleteStep(ctx, apprenticeId, event)
}

func (im *InvitationManagerFacade) ListApprentices(ctx context.Context, masterId uint64, depth int) ([]RelationWithRewardAmount, error) {
	m, err := im.getManager(ctx)
	if err != nil {
		return nil, err
	}
	return m.ListApprentices(ctx, masterId, depth)
}

func (im *InvitationManagerFacade) GetToken(ctx context.Context, id uint) string {
	m, _ := im.getManager(ctx)
	return m.GetToken(ctx, id)
}

func (im *InvitationManagerFacade) GetUrl(ctx context.Context, claim *jwt2.Claim) string {
	m, _ := im.getManager(ctx)
	return m.GetUrl(ctx, claim)
}

func (im *InvitationManagerFacade) getManager(ctx context.Context) (*InvitationManager, error) {
	tenant := config.GetTenant(ctx)
	conf, err := im.DynConf.Tenant(tenant)
	if err != nil {
		return nil, errors.Wrap(err, "no configuration found for invitation tenant")
	}
	return im.Factory.NewManager(conf.Cut("share")), nil
}
