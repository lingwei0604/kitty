package listener

import (
	"context"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/event"
	kitty "github.com/lingwei0604/kitty/proto"
	shareevent "github.com/lingwei0604/kitty/share/event"
)

type InvitationCodeBus interface {
	Emit(ctx context.Context, info contract.Marshaller) error
}

type InvitationCodeAdded struct {
	Bus InvitationCodeBus
}

func (i InvitationCodeAdded) Listen() []contract.Event {
	return event.Of(shareevent.InvitationCodeAdded{})
}

func (i InvitationCodeAdded) Process(event contract.Event) error {
	var info *kitty.InvitationInfo

	data := event.Data().(shareevent.InvitationCodeAdded)

	info = &kitty.InvitationInfo{
		InviteeId:   data.InviteeId,
		InviterId:   data.InviterId,
		DateTime:    data.DateTime.Format("2006-01-02 15:04:05"),
		PackageName: data.PackageName,
		Channel:     data.Channel,
		Ipv4:        data.Ipv4,
	}

	_ = i.Bus.Emit(event.Context(), info)
	return nil
}
