package wechat

import (
	"strings"
	"sync"

	"github.com/lingwei0604/kitty/pkg/contract"
)

type WechaterFactory struct {
	mutex sync.Mutex
	cache map[string]Wechater
	conf  contract.ConfigReader
	doer  contract.HttpDoer
}

func NewWechaterFactory(conf contract.ConfigReader, doer contract.HttpDoer) *WechaterFactory {
	return &WechaterFactory{conf: conf, doer: doer, cache: make(map[string]Wechater)}
}

func (t *WechaterFactory) getConfig(conf contract.ConfigReader, name string) *WechatConfig {
	conf = conf.Cut(name)
	return &WechatConfig{
		WechatAccessTokenUrl: conf.String("wechat.wechatAccessTokenUrl"),
		WeChatGetUserInfoUrl: conf.String("wechat.wechatGetUserInfoUrl"),
		AppId:                conf.String("wechat.appId"),
		AppSecret:            conf.String("wechat.appSecret"),
		Client:               t.doer,
	}
}

func (t *WechaterFactory) GetTransport(name string) Wechater {
	return t.GetTransportWithConf(name, t.conf)
}

func (t *WechaterFactory) GetTransportWithConf(name string, conf contract.ConfigReader) Wechater {
	name = strings.ReplaceAll(name, ".", "_")
	if name == "" {
		name = "default"
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if tsp, ok := t.cache[name]; ok {
		return tsp
	}
	// Currently we only have one kind of sender.
	t.cache[name] = NewTransport(t.getConfig(conf, name))
	return t.cache[name]
}

func (t *WechaterFactory) GetTransportByConf(conf contract.ConfigReader) Wechater {
	// Currently we only have one kind of sender.
	return NewTransport(&WechatConfig{
		WechatAccessTokenUrl: conf.String("wechat.wechatAccessTokenUrl"),
		WeChatGetUserInfoUrl: conf.String("wechat.wechatGetUserInfoUrl"),
		AppId:                conf.String("wechat.appId"),
		AppSecret:            conf.String("wechat.appSecret"),
		Client:               t.doer,
	})
}
