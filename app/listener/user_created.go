package listener

import (
	"context"

	appevent "git.yingzhongshare.com/mkt/kitty/app/event"
	"git.yingzhongshare.com/mkt/kitty/pkg/config"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/event"
)

type EventBus interface {
	Emit(ctx context.Context, event string, tenant *config.Tenant) error
}

type UserCreated struct {
	Bus EventBus
}

func (u UserCreated) Listen() []contract.Event {
	return event.Of(appevent.UserCreated{})
}

func (u UserCreated) Process(event contract.Event) error {
	data := event.Data().(appevent.UserCreated)
	claim := config.Tenant{
		PackageName: data.PackageName,
		UserId:      data.Id,
		Suuid:       data.Suuid,
		Channel:     data.Channel,
		VersionCode: data.VersionCode,
		Os:          uint8(data.Os),
		Idfa:        data.Idfa,
		Imei:        data.Imei,
		Mac:         data.Mac,
		Oaid:        data.Oaid,
		Ip:          data.Ip,
		AndroidId:   data.AndroidId,
	}
	_ = u.Bus.Emit(event.Context(), "new_user", &claim)
	return nil
}
