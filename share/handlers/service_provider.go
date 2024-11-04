package handlers

import (
	"github.com/lingwei0604/kitty/pkg/contract"
	code "github.com/lingwei0604/kitty/pkg/invitecode"
	pb "github.com/lingwei0604/kitty/proto"
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
