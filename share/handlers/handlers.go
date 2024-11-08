//go:generate mockery --name=InvitationManager
//go:generate mockery --name=UserRepository
package handlers

import (
	"context"
	"time"

	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/event"
	shareevent "github.com/lingwei0604/kitty/share/event"

	"github.com/lingwei0604/kitty/app/entity"
	"github.com/lingwei0604/kitty/app/msg"
	token "github.com/lingwei0604/kitty/pkg/invitecode"
	"github.com/lingwei0604/kitty/pkg/kerr"
	kittyjwt "github.com/lingwei0604/kitty/pkg/kjwt"
	pb "github.com/lingwei0604/kitty/proto"
	"github.com/lingwei0604/kitty/share/internal"
	"github.com/pkg/errors"
)

var ErrReenteringInviteCode = errors.New("不能重复填写邀请码")

type shareService struct {
	manager    InvitationManager
	ur         UserRepository
	dispatcher contract.Dispatcher
	tokenizer  internal.EncodeDecoder
}

type InvitationManager interface {
	AddToken(ctx context.Context, apprentice, master *entity.User) error
	ClaimReward(ctx context.Context, masterId uint64, apprenticeId uint64) error
	CompleteStep(ctx context.Context, apprenticeId uint64, event internal.ReceivedEvent) error
	ListApprentices(ctx context.Context, masterId uint64, depth int) ([]internal.RelationWithRewardAmount, error)
	ListMaster(ctx context.Context, apprenticeId uint64) (master *entity.User, grandMaster *entity.User, err error)
	GetToken(ctx context.Context, id uint) string
	GetUrl(ctx context.Context, claim *kittyjwt.Claim) string
}

type UserRepository interface {
	UpdateCallback(ctx context.Context, id uint, f func(user *entity.User) error) (err error)
	Get(ctx context.Context, id uint) (*entity.User, error)
}

func (s shareService) AddInvitationCode(ctx context.Context, in *pb.ShareAddInvitationRequest) (*pb.ShareGenericReply, error) {
	var err error

	claim := kittyjwt.ClaimFromContext(ctx)

	inviterId, err := s.tokenizer.Decode(in.GetInviteCode())

	if err != nil {
		return nil, kerr.FailedPreconditionErr(err, msg.InvalidInviteCode)
	}

	master, err := s.ur.Get(ctx, inviterId)
	if err != nil {
		return nil, kerr.FailedPreconditionErr(err, msg.InvalidInviteCode)
	}

	err = s.ur.UpdateCallback(ctx, uint(claim.UserId), func(user *entity.User) error {
		if user.InviteCode != "" {
			return ErrReenteringInviteCode
		}

		err := s.manager.AddToken(ctx, user, master)
		if err != nil {
			return errors.Wrap(err, msg.InvalidInviteCode)
		}
		user.InviteCode = in.InviteCode
		return nil
	})

	if errors.Is(err, ErrReenteringInviteCode) {
		return nil, kerr.InvalidArgumentErr(err, msg.ReenteringCode)
	}
	if errors.Is(err, entity.ErrRelationArgument) {
		return nil, kerr.InvalidArgumentErr(err, msg.InvalidInviteTarget)
	}
	if errors.Is(err, entity.ErrRelationSequence) {
		return nil, kerr.InvalidArgumentErr(err, msg.InvalidInviteSequence)
	}
	if errors.Is(err, entity.ErrRelationCircled) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorCircledInvitation)
	}
	if errors.Is(err, entity.ErrRelationExist) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorRelationAlreadyExists)
	}
	if errors.Is(err, token.ErrFailedToDecodeToken) {
		return nil, kerr.FailedPreconditionErr(err, msg.InvalidInviteCode)
	}

	// 触发事件
	e := shareevent.InvitationCodeAdded{
		InviteeId:   claim.UserId,
		InviterId:   uint64(inviterId),
		PackageName: claim.PackageName,
		InviteCode:  in.GetInviteCode(),
		DateTime:    time.Now(),
		Channel:     claim.Channel,
	}

	if ip, ok := ctx.Value(contract.IpKey).(string); ok {
		e.Ipv4 = ip
	}

	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, e))

	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) ClaimReward(ctx context.Context, in *pb.ShareClaimRewardRequest) (*pb.ShareGenericReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	err := s.manager.ClaimReward(ctx, claim.UserId, in.ApprenticeId)
	if err != nil {
		if errors.Is(err, entity.ErrOrientationHasNotBeenCompleted) {
			return nil, kerr.FailedPreconditionErr(err, msg.OrientationHasNotBeenCompleted)
		}
		if errors.Is(err, entity.ErrRewardClaimed) {
			return nil, kerr.FailedPreconditionErr(err, msg.RewardClaimed)
		}
		if errors.Is(err, internal.ErrFailedXtaskRequest) {
			return nil, kerr.FailedPreconditionErr(err, msg.XTastAbnormally)
		}
		return nil, kerr.InternalErr(err, msg.NoRewardAvailable)
	}
	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) ListFriend(ctx context.Context, in *pb.ShareListFriendRequest) (*pb.ShareListFriendReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	rels, err := s.manager.ListApprentices(ctx, claim.UserId, int(in.Depth))
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var (
		resp          pb.ShareListFriendReply
		countNotReady int32
		countReady    int32
		countClaimed  int32
	)
	resp.Data = new(pb.ShareListFriendData)
	for _, rel := range rels {
		item := &pb.ShareListFriendDataItem{
			Id:       uint64(rel.ApprenticeID),
			UserName: rel.Apprentice.UserName,
			HeadImg:  rel.Apprentice.HeadImg,
			Gender:   pb.Gender(rel.Apprentice.Gender),
			Coin:     int32(rel.Amount),
			Steps:    make(map[string]bool),
			CreateAt: rel.CreatedAt.Unix(),
		}
		item.ClaimStatus = status(&rel, &countNotReady, &countReady, &countClaimed)

		for _, step := range rel.OrientationSteps {
			item.Steps[step.ChineseName] = step.StepCompleted
		}
		resp.Data.Items = append(resp.Data.Items, item)
	}
	resp.Data.CountNotReady = countNotReady
	resp.Data.CountAll = int32(len(rels))
	resp.Data.CountReady = countReady
	resp.Data.CountClaimed = countClaimed
	return &resp, nil
}

func status(item *internal.RelationWithRewardAmount, countNotReady *int32, countReady *int32, countClaimed *int32) pb.ClaimStatus {
	if item.RewardClaimed {
		*countClaimed++
		return pb.ClaimStatus_DONE
	}
	if item.OrientationCompleted {
		*countReady++
		return pb.ClaimStatus_READY
	}
	*countNotReady++
	return pb.ClaimStatus_NOT_READY
}

func (s shareService) InviteByUrl(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataUrlReply, error) {
	url := s.manager.GetUrl(ctx, kittyjwt.ClaimFromContext(ctx))
	var resp = pb.ShareDataUrlReply{
		Code: 0,
		Data: &pb.ShareDataUrlReply_Url{
			Url: url,
		},
	}
	return &resp, nil
}

func (s shareService) InviteByToken(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataTokenReply, error) {
	id := uint(kittyjwt.ClaimFromContext(ctx).UserId)
	code := s.manager.GetToken(ctx, id)
	var resp = pb.ShareDataTokenReply{
		Code: 0,
		Data: &pb.ShareDataTokenReply_Code{
			Code: code,
		},
	}
	return &resp, nil
}

func (s shareService) PushSignEvent(ctx context.Context, in *pb.SignEvent) (*pb.ShareGenericReply, error) {
	var event internal.ReceivedEvent
	event.Id = int(in.Id)
	event.Type = "sign"

	err := s.manager.CompleteStep(ctx, in.UserId, event)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) PushTaskEvent(ctx context.Context, in *pb.TaskEvent) (*pb.ShareGenericReply, error) {

	var event internal.ReceivedEvent
	event.Id = int(in.Id)
	event.Type = "task"

	err := s.manager.CompleteStep(ctx, in.UserId, event)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var resp pb.ShareGenericReply
	return &resp, nil
}

func (s shareService) GetMaster(ctx context.Context, in *pb.ShareGetMasterRequest) (*pb.ShareGetMasterReply, error) {
	var masterInfo, grandMasterInfo *pb.UserInfo
	if in.Id == 0 {
		claim := kittyjwt.ClaimFromContext(ctx)
		in.Id = claim.UserId
	}

	master, grandMaster, err := s.manager.ListMaster(ctx, in.Id)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	if master != nil {
		masterInfo = &pb.UserInfo{
			Id:           uint64(master.ID),
			UserName:     master.UserName,
			Wechat:       master.WechatOpenId.String,
			HeadImg:      master.HeadImg,
			Gender:       pb.Gender(master.Gender),
			Birthday:     master.Birthday,
			ThirdPartyId: master.ThirdPartyId,
			IsNew:        master.IsNew,
			IsInvited:    master.InviteCode != "",
			CreatedAt:    master.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	if grandMaster != nil {
		grandMasterInfo = &pb.UserInfo{
			Id:           uint64(grandMaster.ID),
			UserName:     grandMaster.UserName,
			Wechat:       grandMaster.WechatOpenId.String,
			HeadImg:      grandMaster.HeadImg,
			Gender:       pb.Gender(grandMaster.Gender),
			Birthday:     grandMaster.Birthday,
			ThirdPartyId: grandMaster.ThirdPartyId,
			IsNew:        grandMaster.IsNew,
			IsInvited:    grandMaster.InviteCode != "",
			CreatedAt:    grandMaster.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &pb.ShareGetMasterReply{
		Code: 0,
		Data: &pb.ShareGetMasterReply_Data{
			Master:      masterInfo,
			GrandMaster: grandMasterInfo,
		},
	}, nil
}
