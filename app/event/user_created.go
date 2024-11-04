package event

import pb "git.yingzhongshare.com/mkt/kitty/proto"

type UserCreated struct {
	*pb.UserInfoDetail
}
