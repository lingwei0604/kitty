package listener

import (
	"context"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/event"
	kitty "git.yingzhongshare.com/mkt/kitty/proto"
	shareevent "git.yingzhongshare.com/mkt/kitty/share/event"
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
