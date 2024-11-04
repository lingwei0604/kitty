//go:generate mockery --name=CodeRepository
//go:generate mockery --name=UserRepository
//go:generate mockery --name=FileRepository
//go:generate mockery --name=PreAllocator
//go:generate mockery --name=UserCache

package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lingwei0604/kitty/app/entity"
	appevent "github.com/lingwei0604/kitty/app/event"
	"github.com/lingwei0604/kitty/app/msg"
	"github.com/lingwei0604/kitty/app/repository"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/event"
	code "github.com/lingwei0604/kitty/pkg/invitecode"
	"github.com/lingwei0604/kitty/pkg/kerr"
	kittyjwt "github.com/lingwei0604/kitty/pkg/kjwt"
	"github.com/lingwei0604/kitty/pkg/wechat"
	pb "github.com/lingwei0604/kitty/proto"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type appService struct {
	conf         contract.ConfigReader
	logger       log.Logger
	ur           UserRepository
	cr           CodeRepository
	fr           FileRepository
	sender       contract.SmsSender
	wechat       wechat.Wechater
	dispatcher   contract.Dispatcher
	preAllocator PreAllocator
	userCache    UserCache
}

type tokenParam struct {
	userId                                                                 uint64
	suuid, channel, versionCode, wechat, mobile, packageName, thirdPartyId string
	registerTime                                                           int64
}

type CodeRepository interface {
	CheckCode(ctx context.Context, mobile, code string) (bool, error)
	AddCode(ctx context.Context, mobile string) (code string, err error)
	DeleteCode(ctx context.Context, mobile string) (err error)
}

type UserRepository interface {
	GetFromWechat(ctx context.Context, packageName, wechat string, device *entity.Device, wechatUser entity.User) (user *entity.User, err error)
	GetFromMobile(ctx context.Context, packageName, mobile string, device *entity.Device) (user *entity.User, err error)
	GetFromDevice(ctx context.Context, packageName, suuid string, device *entity.Device) (user *entity.User, err error)
	GetRecentUserByDevice(ctx context.Context, packageName string, device *entity.Device) (*entity.User, error)
	GetByDevice(ctx context.Context, packageName string, device entity.Device) (user []entity.User, err error)
	GetIDAndCreatedAtBySUUID(ctx context.Context, packageName, suuid string) (id uint, createdAt time.Time, err error)
	Update(ctx context.Context, id uint, user entity.User) (newUser *entity.User, err error)
	Get(ctx context.Context, id uint) (user *entity.User, err error)
	GetAll(ctx context.Context, where ...clause.Expression) (user []entity.User, err error)
	Count(ctx context.Context, where ...clause.Expression) (total int64, err error)
	Save(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	DeleteDevices(ctx context.Context, uid uint) ([]entity.Device, error)
	DeleteDevicesByOaid(context.Context, string) ([]entity.Device, error)
	DeleteDevicesByIdfa(context.Context, string) ([]entity.Device, error)
	DeleteDevicesByAndroid(context.Context, string) ([]entity.Device, error)
	GetByWechat(ctx context.Context, packageName, wechat string) (user *entity.User, err error)
}

type PreAllocator interface {
	Preallocate(ctx context.Context, packageName, suuid string) (uint, error)
	ClearPreallocate(ctx context.Context, packageName, suuid string) error

	SetIDByDevice(ctx context.Context, packageName string, device *entity.Device, ID uint) error
	GetIDByDevice(ctx context.Context, packageName string, device *entity.Device) (uint, error)
	SetRegisterTimeById(ctx context.Context, ID uint, registerTime int64) error
	GetRegisterTimeById(ctx context.Context, ID uint) (int64, error)
	DeleteCacheByDevice(ctx context.Context, packageName string, device *entity.Device)
}

type UserCache interface {
	CacheID(ctx context.Context, id uint) error
	CacheBindAd(ctx context.Context, cb *repository.BindAdCallback) error
	GetBindAd(ctx context.Context, id uint) (*repository.BindAdCallback, error)
}

type FileRepository interface {
	UploadFromUrl(ctx context.Context, url string) (newUrl string, err error)
}

func (s appService) Login(ctx context.Context, in *pb.UserLoginRequest) (resp *pb.UserInfoReply, err error) {
	var (
		u      *entity.User
		device *entity.Device
	)

	defer func() {
		reqB, _ := json.Marshal(in)
		respB, _ := json.Marshal(resp)
		level.Warn(s.logger).Log("action", "login", "loginReq", string(reqB), "response", string(respB))
	}()

	device = &entity.Device{
		Os:        uint8(in.Device.Os),
		Imei:      in.Device.Imei,
		Idfa:      in.Device.Idfa,
		Suuid:     in.Device.Suuid,
		Oaid:      in.Device.Oaid,
		Mac:       in.Device.Mac,
		AndroidId: in.Device.AndroidId,
		SMID:      in.Device.Smid,
	}
	if ip, ok := ctx.Value(contract.IpKey).(string); ok {
		device.IP = ip
	}
	u, err = s.loginFrom(ctx, in, device)
	if err != nil {
		return nil, err
	}
	// 预注册Key存在7天,故7天内登录注册过，便可删除预注册Key了
	if time.Now().Sub(u.CreatedAt) < time.Hour*24*7 {
		s.preAllocator.ClearPreallocate(ctx, in.PackageName, device.Suuid)
	}
	s.userCache.CacheID(ctx, u.ID)

	// 再存一些信息
	if err := s.addChannelAndVersionInfo(ctx, in, u); err != nil {
		s.warn(err)
	}

	// Create jwt token
	tokenString, err := s.getToken(&tokenParam{uint64(u.ID), u.CommonSUUID, u.Channel, u.VersionCode, u.WechatOpenId.String, u.Mobile.String, u.PackageName, u.ThirdPartyId, u.CreatedAt.Unix()})
	if err != nil {
		s.warn(err)
	}

	// 触发事件
	var detail = s.toDetail(u)
	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserChanged{UserInfoDetail: detail}))
	if u.IsNew {
		_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserCreated{UserInfoDetail: detail}))
	}

	// 拼装返回结果
	resp = s.toReply(u)
	resp.Data.Token = tokenString

	return resp, nil
}

func (s appService) BindWechat(ctx context.Context, in *pb.BindWechatRequest) (resp *pb.UserInfoReply, err error) {
	var (
		u      *entity.User
		device *entity.Device
	)

	defer func() {
		reqB, _ := json.Marshal(in)
		respB, _ := json.Marshal(resp)
		level.Warn(s.logger).Log("action", "login", "loginReq", string(reqB), "response", string(respB))
	}()

	device = &entity.Device{
		Os:        uint8(in.Device.Os),
		Imei:      in.Device.Imei,
		Idfa:      in.Device.Idfa,
		Suuid:     in.Device.Suuid,
		Oaid:      in.Device.Oaid,
		Mac:       in.Device.Mac,
		AndroidId: in.Device.AndroidId,
		SMID:      in.Device.Smid,
	}
	if ip, ok := ctx.Value(contract.IpKey).(string); ok {
		device.IP = ip
	}
	// 查询 openid 对应了 user_info
	u, err = s.bindWechat(ctx, in, device)
	if err != nil {
		return nil, err
	}

	// Create jwt token
	tokenString, err := s.getToken(&tokenParam{uint64(u.ID), u.CommonSUUID, u.Channel, u.VersionCode, u.WechatOpenId.String, u.Mobile.String, u.PackageName, u.ThirdPartyId, u.CreatedAt.Unix()})
	if err != nil {
		s.warn(err)
	}

	// 触发事件
	var detail = s.toDetail(u)
	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserChanged{UserInfoDetail: detail}))
	if u.IsNew {
		_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserCreated{UserInfoDetail: detail}))
	}

	// 拼装返回结果
	resp = s.toReply(u)
	resp.Data.Token = tokenString

	return resp, nil
}

func (s appService) GetCode(ctx context.Context, in *pb.GetCodeRequest) (*pb.GenericReply, error) {
	code, err := s.cr.AddCode(ctx, in.Mobile)
	if err == repository.ErrTooFrequent {
		return nil, kerr.ResourceExhaustedErr(err, msg.ErrorTooFrequent)
	}
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorGetCode)
	}
	err = s.sender.Send(ctx, in.Mobile, code)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorSendCode)
	}
	var resp = pb.GenericReply{
		Code: 0,
	}
	return &resp, nil
}

func (s appService) GetInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	if in.Id == 0 {
		claim := kittyjwt.ClaimFromContext(ctx)
		in.Id = claim.UserId
	}
	u, err := s.ur.Get(ctx, uint(in.Id))
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, kerr.NotFoundErr(err, msg.ErrorRecordNotFound)
	}
	if err != nil {
		return nil, dbErr(err)
	}
	var resp = s.toReply(u)

	if !in.Taobao {
		resp.Data.TaobaoExtra = nil
	}

	if !in.Wechat {
		resp.Data.WechatExtra = nil
	}

	return resp, nil
}

func (s appService) Refresh(ctx context.Context, in *pb.UserRefreshRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	device := &entity.Device{
		Os:        uint8(in.Device.Os),
		Imei:      in.Device.Imei,
		Idfa:      in.Device.Idfa,
		Suuid:     in.Device.Suuid,
		Oaid:      in.Device.Oaid,
		Mac:       in.Device.Mac,
		AndroidId: in.Device.AndroidId,
		SMID:      in.Device.Smid,
	}
	if ip, ok := ctx.Value(contract.IpKey).(string); ok {
		device.IP = ip
	}

	u, err := s.ur.Get(ctx, uint(claim.UserId))
	if err != nil {
		return nil, dbErr(err)
	}

	u.Channel = in.Channel
	u.VersionCode = in.VersionCode
	u.CommonSMID = device.SMID
	if device.Suuid != "" {
		u.CommonSUUID = device.Suuid
	}
	u.AddNewDevice(device)

	if err := s.ur.Save(ctx, u); err != nil {
		return nil, errors.Wrap(err, "Refresh Save")
	}

	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserChanged{UserInfoDetail: s.toDetail(u)}))
	reply := s.toReply(u)
	reply.Data.Token, err = s.getToken(&tokenParam{
		uint64(u.ID),
		u.CommonSUUID,
		u.Channel,
		u.VersionCode,
		u.WechatOpenId.String,
		u.Mobile.String,
		u.PackageName,
		u.ThirdPartyId,
		u.CreatedAt.Unix(),
	})

	if err != nil {
		err = kerr.InternalErr(err, msg.ErrorJwtFailure)
	}
	return reply, nil
}

func (s appService) GetInfoBatch(ctx context.Context, in *pb.UserInfoBatchRequest) (*pb.UserInfoBatchReply, error) {
	var expressions []clause.Expression
	if len(in.Id) > 0 {
		var ids []interface{}
		for _, v := range in.Id {
			ids = append(ids, uint(v))
		}
		expressions = append(expressions, clause.IN{
			Column: "id",
			Values: ids,
		})
	}
	if len(in.InviteCode) > 0 {
		var ids []interface{}
		for _, v := range in.InviteCode {
			t := code.NewTokenizer(s.conf.String("salt"))
			id, _ := t.Decode(v)
			ids = append(ids, id)
		}
		expressions = append(expressions, clause.IN{
			Column: "id",
			Values: ids,
		})
	}
	if len(in.PackageName) > 0 {
		expressions = append(expressions, clause.Eq{
			Column: "package_name",
			Value:  in.PackageName,
		})
	}
	if in.After != 0 {
		expressions = append(expressions, clause.Gt{
			Column: "created_at",
			Value:  time.Unix(in.After, 0),
		})
	}
	if in.Before != 0 {
		expressions = append(expressions, clause.Lt{
			Column: "created_at",
			Value:  time.Unix(in.After, 0),
		})
	}
	if len(in.Name) != 0 {
		expressions = append(expressions, clause.Like{
			Column: "user_name",
			Value:  "%" + in.Name + "%",
		})
	}
	if len(in.Mobile) != 0 {
		expressions = append(expressions, clause.Eq{
			Column: "mobile",
			Value:  in.Mobile,
		})
	}

	c := clause.Where{
		Exprs: expressions,
	}

	count, err := s.ur.Count(ctx, c)
	if err != nil {
		return nil, dbErr(err)
	}
	if in.PerPage <= 0 {
		in.PerPage = 20
	}
	if in.Page <= 0 {
		in.Page = 1
	}
	limit := clause.Limit{
		Limit:  int(in.PerPage),
		Offset: int((in.Page - 1) * in.PerPage),
	}

	users, err := s.ur.GetAll(ctx, c, limit)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, kerr.NotFoundErr(err, msg.ErrorRecordNotFound)
	}
	if err != nil {
		return nil, dbErr(err)
	}
	var resp = pb.UserInfoBatchReply{
		Code: 0,
		Data: []*pb.UserInfoDetail{},
	}

	for _, v := range users {
		tmp := s.toDetail(&v)
		resp.Data = append(resp.Data, tmp)
	}
	resp.Count = count
	return &resp, nil
}

func (s appService) UpdateInfo(ctx context.Context, in *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	u, err := s.ur.Update(ctx, uint(claim.UserId), entity.User{
		UserName:     in.UserName,
		HeadImg:      in.HeadImg,
		Gender:       int(in.Gender),
		Birthday:     in.Birthday,
		ThirdPartyId: in.ThirdPartyId,
	})
	if err != nil {
		return nil, dbErr(err)
	}
	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserChanged{UserInfoDetail: s.toDetail(u)}))
	var resp = s.toReply(u)
	return resp, nil

}

func (s appService) SoftDelete(ctx context.Context, in *pb.UserSoftDeleteRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	if in.Id != 0 && in.Id != claim.UserId {
		// 删除别人的账号需要管理员权限
		if claim.Audience != "admin" {
			return nil, kerr.UnauthenticatedErr(errors.New("action requires admin privilege"), msg.AdminOnly)
		}
	}
	if in.Id == 0 {
		in.Id = claim.UserId
	}
	u, err := s.unbindId(ctx, &pb.UserUnbindRequest{
		Mobile: true,
		Wechat: true,
		Taobao: true,
		Suuid:  true,
	}, in.Id)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, kerr.NotFoundErr(err, msg.AlreadyDeleted)
	}
	if err != nil {
		return nil, err
	}
	err = s.ur.Delete(ctx, uint(in.Id))
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, dbErr(err)
	}
	u.Data.IsDeleted = true
	return u, nil
}

func (s appService) Bind(ctx context.Context, in *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)

	var (
		toUpdate entity.User
		err      error
	)

	// 兼容旧接口
	if len(in.OpenId) > 0 {
		in.WechatExtra = &pb.WechatExtra{OpenId: in.OpenId}
	}

	// 绑定手机号
	if len(in.Mobile) > 0 && len(in.Code) > 0 {
		if ok, err := s.verify(ctx, in.Mobile, in.Code); err != nil {
			return nil, dbErr(err)
		} else if !ok {
			return nil, kerr.UnauthenticatedErr(errors.Errorf("cannot verify %s with %s", in.Mobile, in.Code), msg.ErrorMobileCode)
		}
		toUpdate.Mobile = ns(in.Mobile)
	}

	// 绑定微信号
	if len(in.Wechat) > 0 {
		var wechatExtra *pb.WechatExtra
		wechatExtra, err = s.getWechatInfo(ctx, in.Wechat)
		if err != nil {
			return nil, kerr.UnauthenticatedErr(err, msg.ErrorWechatFailure)
		}
		wechatExtraBytes, err := wechatExtra.Marshal()
		if err != nil {
			return nil, kerr.InternalErr(err, msg.ErrorLogin)
		}
		toUpdate.WechatOpenId = ns(wechatExtra.OpenId)
		toUpdate.WechatUnionId = ns(wechatExtra.Unionid)
		toUpdate.WechatExtra = wechatExtraBytes
		if in.MergeInfo {
			toUpdate.UserName = wechatExtra.NickName
			toUpdate.HeadImg, _ = s.fr.UploadFromUrl(ctx, wechatExtra.Headimgurl)
			toUpdate.Gender = int(wechatExtra.Sex)
		}
	}

	// 绑定淘宝openId
	if in.TaobaoExtra != nil && len(in.TaobaoExtra.OpenId) > 0 {
		taobaoExtraBytes, err := in.TaobaoExtra.Marshal()
		if err != nil {
			return nil, kerr.InternalErr(err, msg.ErrorCorruptedData)
		}
		toUpdate.TaobaoOpenId = ns(in.TaobaoExtra.OpenId)
		toUpdate.TaobaoExtra = taobaoExtraBytes
		if in.MergeInfo {
			toUpdate.UserName = in.TaobaoExtra.Nick
			toUpdate.HeadImg, _ = s.fr.UploadFromUrl(ctx, in.TaobaoExtra.AvatarUrl)
		}
	}

	// 绑定微信openId
	if in.WechatExtra != nil && len(in.WechatExtra.OpenId) > 0 {
		wechatExtraBytes, err := in.WechatExtra.Marshal()
		if err != nil {
			return nil, kerr.InternalErr(err, msg.ErrorCorruptedData)
		}
		toUpdate.WechatOpenId = ns(in.WechatExtra.OpenId)
		toUpdate.WechatExtra = wechatExtraBytes
	}

	// 更新用户
	newUser, err := s.ur.Update(ctx, uint(claim.UserId), toUpdate)
	if errors.Is(err, repository.ErrAlreadyBind) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorAlreadyBind)
	}
	if err != nil {
		return nil, dbErr(err)
	}

	// 获取Token
	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserChanged{UserInfoDetail: s.toDetail(newUser)}))
	reply := s.toReply(newUser)
	reply.Data.Token, err = s.getToken(&tokenParam{
		uint64(newUser.ID),
		newUser.CommonSUUID,
		newUser.Channel,
		newUser.VersionCode,
		newUser.WechatOpenId.String,
		newUser.Mobile.String,
		newUser.PackageName,
		newUser.ThirdPartyId,
		newUser.CreatedAt.Unix(),
	})
	if err != nil {
		err = kerr.InternalErr(err, msg.ErrorJwtFailure)
	}

	return reply, err
}

func (s appService) Unbind(ctx context.Context, in *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	if claim.Audience != "admin" {
		return nil, kerr.UnauthenticatedErr(errors.New("action requires admin privilege"), msg.AdminOnly)
	}
	if in.UserId == 0 {
		in.UserId = int64(claim.UserId)
	}
	if in.Oaid != "" {
		_, err := s.ur.DeleteDevicesByOaid(ctx, in.Oaid)
		if err != nil {
			return nil, kerr.InternalErr(errors.New("delete device by oaid error"), msg.ErrorDatabaseFailure)
		}
	}
	if in.Idfa != "" {
		_, err := s.ur.DeleteDevicesByIdfa(ctx, in.Idfa)
		if err != nil {
			return nil, kerr.InternalErr(errors.New("delete device by idfa error"), msg.ErrorDatabaseFailure)
		}
	}
	if in.Android != "" {
		_, err := s.ur.DeleteDevicesByAndroid(ctx, in.Android)
		if err != nil {
			return nil, kerr.InternalErr(errors.New("delete device by android error"), msg.ErrorDatabaseFailure)
		}
	}
	return s.unbindId(ctx, in, uint64(in.UserId))
}

func (s appService) unbindId(ctx context.Context, in *pb.UserUnbindRequest, id uint64) (*pb.UserInfoReply, error) {
	user, err := s.ur.Get(ctx, uint(id))
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, kerr.NotFoundErr(err, msg.ErrorRecordNotFound)
	}
	if err != nil {
		return nil, err
	}
	if in.Mobile {
		user.Mobile = sql.NullString{}
	}
	if in.Wechat {
		user.WechatUnionId = sql.NullString{}
		user.WechatOpenId = sql.NullString{}
		user.WechatExtra = nil
	}
	if in.Taobao {
		user.TaobaoOpenId = sql.NullString{}
		user.TaobaoExtra = nil
	}
	var suuid = user.CommonSUUID
	if in.Suuid {
		user.CommonSUUID = ""
	}

	err = s.ur.Save(ctx, user)
	if err != nil {
		return nil, dbErr(err)
	}
	if in.Suuid && suuid != "" {
		s.preAllocator.ClearPreallocate(ctx, user.PackageName, suuid)
	}
	devices, err := s.ur.DeleteDevices(ctx, user.ID)
	if err != nil {
		level.Error(s.logger).Log("err", err)
	} else {
		for _, d := range devices {
			s.preAllocator.DeleteCacheByDevice(ctx, user.PackageName, &d)
		}
	}
	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserChanged{UserInfoDetail: s.toDetail(user)}))
	var resp = s.toReply(user)
	return resp, nil
}

func (s appService) getToken(param *tokenParam) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		kittyjwt.NewClaim(
			param.userId,
			s.conf.String("name"),
			param.suuid, param.channel, param.versionCode, param.wechat, param.mobile, param.packageName,
			param.registerTime,
			time.Hour*24*30,
		),
	)
	token.Header["kid"] = s.conf.String("security.kid")
	return token.SignedString([]byte(s.conf.String("security.key")))
}
func (s appService) debug(err error) {
	if err != nil {
		level.Debug(s.logger).Log("err", err)
	}
}

func (s appService) error(err error) {
	if err != nil {
		level.Error(s.logger).Log("err", err)
	}
}

func (s appService) warn(err error) {
	if err != nil {
		level.Warn(s.logger).Log("err", err)
	}
}

func (s appService) getWechatInfo(ctx context.Context, wechat string) (*pb.WechatExtra, error) {
	wxRes, err := s.wechat.GetLoginResponse(ctx, wechat)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	if wxRes.Openid == "" {
		return nil, errors.New(msg.ErrorMissingOpenid)
	}
	wxInfo, err := s.wechat.GetUserInfoResult(ctx, wxRes)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	// side effect: save extra wechat info
	infoPb := &pb.WechatExtra{
		AccessToken:  wxRes.AccessToken,
		ExpiresIn:    wxRes.ExpiresIn,
		RefreshToken: wxRes.RefreshToken,
		OpenId:       wxRes.Openid,
		Scope:        wxRes.Scope,
		NickName:     wxInfo.Nickname,
		Sex:          int32(wxInfo.Sex),
		Province:     wxInfo.Province,
		City:         wxInfo.City,
		Country:      wxInfo.Country,
		Headimgurl:   wxInfo.Headimgurl,
		Privilege:    wxInfo.Privilege,
		Unionid:      wxInfo.Unionid,
	}
	return infoPb, nil
}

func (s appService) handleWechatLogin(ctx context.Context, packageName, wechat string, device *entity.Device) (*entity.User, error) {
	wxInfo, err := s.getWechatInfo(ctx, wechat)
	if err != nil {
		return nil, kerr.UnauthenticatedErr(err, msg.ErrorWechatFailure)
	}
	extra, err := wxInfo.Marshal()
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorCorruptedData)
	}

	wechatUser := entity.User{
		UserName:      wxInfo.NickName,
		HeadImg:       wxInfo.Headimgurl,
		WechatOpenId:  ns(wxInfo.OpenId),
		WechatUnionId: ns(wxInfo.Unionid),
		WechatExtra:   extra,
	}

	u, err := s.ur.GetFromWechat(ctx, packageName, wxInfo.OpenId, device, wechatUser)
	if err != nil {
		return nil, dbErr(err)
	}
	level.Info(s.logger).Log("msg", fmt.Sprintf(msg.WxSuccess, u.ID), "suuid", device.Suuid, "userId", u.ID, "packageName", packageName)
	return u, nil
}

func (s appService) handleMobileLogin(ctx context.Context, packageName, mobile, code string, device *entity.Device) (*entity.User, error) {
	if len(code) == 0 {
		return nil, kerr.InvalidArgumentErr(errors.New("code cannot be 0"), msg.InvalidParams)
	}
	if ok, err := s.verify(ctx, mobile, code); err != nil {
		return nil, err
	} else if !ok {
		return nil, kerr.UnauthenticatedErr(errors.Errorf("cannot verify %s with %s", mobile, code), msg.ErrorMobileCode)
	}
	u, err := s.ur.GetFromMobile(ctx, packageName, mobile, device)
	if err != nil {
		return nil, dbErr(err)
	}
	level.Info(s.logger).Log("msg", fmt.Sprintf(msg.MobileSuccess, u.ID), "suuid", device.Suuid, "userId", u.ID, "packageName", packageName)
	return u, nil
}

func (s appService) handleDeviceLogin(ctx context.Context, packageName, suuid string, device *entity.Device) (*entity.User, error) {
	u, err := s.ur.GetFromDevice(ctx, packageName, suuid, device)
	if err != nil {
		return nil, dbErr(err)
	}
	level.Info(s.logger).Log("msg", fmt.Sprintf(msg.DeviceSuccess, u.ID), "suuid", device.Suuid, "userId", u.ID, "packageName", packageName)
	return u, nil
}

func (s appService) loginFrom(ctx context.Context, in *pb.UserLoginRequest, device *entity.Device) (*entity.User, error) {

	if len(in.Mobile) != 0 {
		return s.handleMobileLogin(ctx, in.PackageName, in.Mobile, in.Code, device)
	}

	if len(in.Wechat) != 0 {
		return s.handleWechatLogin(ctx, in.PackageName, in.Wechat, device)
	}

	return s.handleDeviceLogin(ctx, in.PackageName, device.Suuid, device)
}

func (s appService) bindWechat(ctx context.Context, in *pb.BindWechatRequest, device *entity.Device) (*entity.User, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	// 查询 code 对应 open_id
	wxInfo, err := s.getWechatInfo(ctx, in.Wechat)
	if err != nil {
		return nil, kerr.UnauthenticatedErr(err, msg.ErrorWechatFailure)
	}

	// 查询 openid 对应了 user_info
	userInfo, err := s.ur.GetByWechat(ctx, in.PackageName, wxInfo.OpenId)
	if err != nil && err != repository.ErrRecordNotFound {
		return nil, kerr.InternalErr(err, msg.ErrorRecordNotFound)
	}
	// 存在返回用户
	if userInfo != nil {
		return userInfo, nil
	}

	// 绑定微信号
	var toUpdate entity.User
	extra, err := wxInfo.Marshal()
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorCorruptedData)
	}
	toUpdate = entity.User{
		UserName:      wxInfo.NickName,
		HeadImg:       wxInfo.Headimgurl,
		WechatOpenId:  ns(wxInfo.OpenId),
		WechatUnionId: ns(wxInfo.Unionid),
		WechatExtra:   extra,
		Gender:        int(wxInfo.Sex),
	}
	// 更新用户
	newUser, err := s.ur.Update(ctx, uint(claim.UserId), toUpdate)
	if errors.Is(err, repository.ErrAlreadyBind) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorAlreadyBind)
	}
	return newUser, nil
}

func (s appService) addChannelAndVersionInfo(ctx context.Context, in *pb.UserLoginRequest, u *entity.User) error {
	var (
		err      error
		hasExtra bool
	)
	if in.ThirdPartyId != "" && in.ThirdPartyId != u.ThirdPartyId {
		u.ThirdPartyId = in.ThirdPartyId
		hasExtra = true
	}

	if u.Channel != in.Channel {
		u.Channel = in.Channel
		hasExtra = true
	}
	if in.VersionCode != u.VersionCode {
		u.VersionCode = in.VersionCode
		hasExtra = true
	}
	if u.HeadImg == "" {
		u.HeadImg = "https://ad-static-xg.tagtic.cn/ad-material/file/0b8f18e1e666474291174ba316cccb51.png"
	}

	if u.HeadImg == "http://ad-static-xg.tagtic.cn/ad-material/file/0b8f18e1e666474291174ba316cccb51.png" {
		u.HeadImg = "https://ad-static-xg.tagtic.cn/ad-material/file/0b8f18e1e666474291174ba316cccb51.png"
	}

	// 因为bind-ad与用户注册到kitty_user表中存在时间差，故在这里进行补充
	if u.CID.String == "" {
		if cb, err := s.userCache.GetBindAd(ctx, u.ID); err != nil {
			if err != redis.Nil {
				level.Warn(s.logger).Log("msg", "get bindad callback error", "err", err.Error())
			}
		} else if cb != nil {
			hasExtra = true
			u.CampaignID = ns(cb.CampaignId)
			u.AID = ns(cb.Aid)
			u.CID = ns(cb.Cid)
			u.UnionSite = cb.UnionSite
		}
	}

	if hasExtra {
		err = s.ur.Save(ctx, u)
		if err != nil {
			return dbErr(err)
		}
	}
	return nil
}

func (s appService) verify(ctx context.Context, mobile string, code string) (bool, error) {
	result, err := s.cr.CheckCode(ctx, mobile, code)
	if err != nil {
		return false, dbErr(err)
	}
	if result {
		err = s.cr.DeleteCode(ctx, mobile)
		s.warn(err)
	}
	return result, nil
}

func dbErr(err error) kerr.ServerError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return kerr.InternalErr(err, msg.UserNotFound)
	}
	return kerr.InternalErr(err, msg.ErrorDatabaseFailure)
}

func ns(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func redact(mobile string) string {
	if len(mobile) >= 11 {
		mobile = mobile[:3] + "****" + mobile[7:]
	}
	return mobile
}

func (s appService) toReply(user *entity.User) *pb.UserInfoReply {
	var wechatExtra pb.WechatExtra
	if user.WechatOpenId.Valid {
		_ = wechatExtra.Unmarshal(user.WechatExtra)
	}

	var taobaoExtra pb.TaobaoExtra
	if user.TaobaoOpenId.Valid {
		_ = taobaoExtra.Unmarshal(user.TaobaoExtra)
	}
	var tokenizer = code.NewTokenizer(s.conf.String("salt"))
	inviteCode, _ := tokenizer.Encode(user.ID)
	return &pb.UserInfoReply{
		Code:    0,
		Message: "",
		Data: &pb.UserInfo{
			Id:           uint64(user.ID),
			UserName:     user.UserName,
			Wechat:       user.WechatOpenId.String,
			HeadImg:      user.HeadImg,
			Gender:       pb.Gender(user.Gender),
			Birthday:     user.Birthday,
			ThirdPartyId: user.ThirdPartyId,
			Mobile:       redact(user.Mobile.String),
			IsNew:        user.IsNew,
			WechatExtra:  &wechatExtra,
			TaobaoExtra:  &taobaoExtra,
			InviteCode:   inviteCode,
			IsInvited:    user.InviteCode != "",
			CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedTs:    user.CreatedAt.UnixMilli(),
			Smid:         user.CommonSMID,
		},
	}
}

func (s appService) toDetail(user *entity.User) *pb.UserInfoDetail {
	var wechatExtra pb.WechatExtra
	_ = wechatExtra.Unmarshal(user.WechatExtra)
	var taobaoExtra pb.TaobaoExtra
	_ = taobaoExtra.Unmarshal(user.TaobaoExtra)
	var tokenizer = code.NewTokenizer(s.conf.String("salt"))
	inviteCode, _ := tokenizer.Encode(user.ID)
	details := &pb.UserInfoDetail{
		Id:           uint64(user.ID),
		UserName:     user.UserName,
		Wechat:       user.WechatOpenId.String,
		HeadImg:      user.HeadImg,
		Gender:       pb.Gender(user.Gender),
		Birthday:     user.Birthday,
		ThirdPartyId: user.ThirdPartyId,
		Mobile:       redact(user.Mobile.String),
		IsNew:        user.IsNew,
		WechatExtra:  &wechatExtra,
		TaobaoExtra:  &taobaoExtra,
		InviteCode:   inviteCode,
		IsInvited:    user.InviteCode != "",
		Suuid:        user.CommonSUUID,
		Smid:         user.CommonSMID,
		Channel:      user.Channel,
		VersionCode:  user.VersionCode,
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		PackageName:  user.PackageName,
		CampaignId:   user.CampaignID.String,
		Aid:          user.AID.String,
		Cid:          user.CID.String,
		CtaChannel:   user.CtaChannel.String,
	}
	if len(user.Devices) > 0 {
		last := len(user.Devices) - 1
		details.Oaid = user.Devices[last].Oaid
		details.Imei = user.Devices[last].Imei
		details.Idfa = user.Devices[last].Idfa
		details.AndroidId = user.Devices[last].AndroidId
		details.Os = uint32(user.Devices[last].Os)
		details.Ip = user.Devices[last].IP
		details.Mac = user.Devices[last].Mac
	}
	return details
}

func (s appService) BindAd(ctx context.Context, in *pb.UserBindAdRequest) (*pb.GenericReply, error) {
	if in.Id == 0 || in.CampaignId == "" {
		return &pb.GenericReply{}, nil
	}

	u, err := s.ur.Update(ctx, uint(in.Id), entity.User{
		CampaignID: ns(in.CampaignId),
		AID:        ns(in.Aid),
		CID:        ns(in.Cid),
		UnionSite:  in.UnionSite,
		CtaChannel: ns(in.CtaChannel),
	})

	if err != nil {
		if err := s.userCache.CacheBindAd(ctx, &repository.BindAdCallback{
			Id:         in.Id,
			CampaignId: in.CampaignId,
			Cid:        in.Cid,
			Aid:        in.Aid,
			UnionSite:  in.UnionSite,
			CtaChannel: in.CtaChannel,
		}); err != nil {
			level.Warn(s.logger).Log("msg", "cache bindad callback error", "err", err.Error())
		}
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var detail = s.toDetail(u)
	_ = s.dispatcher.Dispatch(event.NewEvent(ctx, appevent.UserChanged{UserInfoDetail: detail}))
	return &pb.GenericReply{}, nil
}

func (s appService) DeviceLookup(ctx context.Context, in *pb.DeviceLookupRequest) (*pb.DeviceLookupReply, error) {
	if claim := kittyjwt.ClaimFromContext(ctx); claim.Audience != "admin" {
		return nil, kerr.UnauthenticatedErr(errors.New("action requires admin privilege"), msg.AdminOnly)
	}

	users, err := s.ur.GetByDevice(ctx, in.PackageName, entity.Device{
		Imei: in.Imei,
		Oaid: in.Oaid,
	})
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	var infos []*pb.UserInfo
	for _, u := range users {
		o := s.toReply(&u)
		infos = append(infos, o.Data)
	}

	return &pb.DeviceLookupReply{
		Code: 0,
		Msg:  "",
		Data: infos,
	}, nil
}

func (s appService) PreRegister(ctx context.Context, in *pb.PreRegisterRequest) (*pb.PreRegisterReply, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var resp pb.PreRegisterReply

	device := &entity.Device{
		Suuid:     in.Suuid,
		Oaid:      in.GetOaid(),
		AndroidId: in.GetAndroidId(),
		Idfa:      in.GetIdfa(),
	}

	defer func() {
		// 缓存设备-uid关系
		s.preAllocator.SetIDByDevice(ctx, in.GetPackageName(), device, uint(resp.UserId))
	}()

	if in.GetUserId() == 0 {
		// 以设备信息在redis缓存中查找
		if id, err := s.preAllocator.GetIDByDevice(ctx, in.GetPackageName(), device); err == nil && id > 0 {
			resp.UserId = uint64(id)
			if reg, _ := s.preAllocator.GetRegisterTimeById(ctx, uint(resp.UserId)); reg > 0 {
				resp.RegisterTime = reg
				return &resp, nil
			}
		}
	}

	// 下方的注册时间均为用户准确的注册时间，
	defer func() {
		s.preAllocator.SetRegisterTimeById(ctx, uint(resp.UserId), resp.RegisterTime)
	}()

	// 传了uid，表示获取用户信息
	if in.GetUserId() != 0 {
		user, err := s.ur.Get(ctx, uint(in.GetUserId()))
		if err == nil && user != nil {
			resp.UserId = uint64(user.ID)
			resp.RegisterTime = user.CreatedAt.UnixMilli()
			return &resp, nil
		}
	}

	// 优先以suuid查询，如果没有，则以设备oaid或idfa查询
	id, createdAt, err := s.ur.GetIDAndCreatedAtBySUUID(ctx, in.PackageName, in.Suuid)
	if err == nil && id != 0 {
		resp.UserId = uint64(id)
		resp.RegisterTime = createdAt.UnixMilli()
		return &resp, nil
	}

	// android oaid值不为空，尝试通过oaid获取用户信息
	user, err := s.ur.GetRecentUserByDevice(ctx, in.GetPackageName(), device)
	if err == nil && user != nil {
		resp.UserId = uint64(user.ID)
		resp.RegisterTime = user.CreatedAt.UnixMilli()
		return &resp, nil
	}

	// 都没有找到，从redis预注册一个uid
	id, err = s.preAllocator.Preallocate(ctx, in.PackageName, in.Suuid)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorDatabaseFailure)
	}
	resp.UserId = uint64(id)
	resp.RegisterTime = time.Now().UnixMilli()
	return &resp, nil
}
