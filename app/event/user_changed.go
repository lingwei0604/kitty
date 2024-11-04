package event

import pb "github.com/lingwei0604/kitty/proto"

type UserChanged struct {
	*pb.UserInfoDetail
}
