package handlers

import (
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	code "git.yingzhongshare.com/mkt/kitty/pkg/invitecode"
	pb "git.yingzhongshare.com/mkt/kitty/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewShareService(
	manager InvitationManager,
	ur UserRepository,
	dispatcher contract.Dispatcher,
	tokenizer *code.Tokenizer,
) *shareService {
	return &shareService{
		manager:    manager,
		ur:         ur,
		dispatcher: dispatcher,
		tokenizer:  tokenizer,
	}
}

func ProvideShareServer(service *shareService) pb.ShareServer {
	return service
}
