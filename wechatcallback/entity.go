package wechatcallback

type WeChatUserInfoChangedPayload struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime   int64  `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Event        string `json:"Event"`
	EventKey     string `json:"EventKey"`
	OpenID       string `json:"OpenID"`
	AppID        string `json:"AppID"`
	RevokeInfo   string `json:"RevokeInfo"`
}
