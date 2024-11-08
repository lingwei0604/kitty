// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.7.0
// source: proto/event_hub.proto

package kitty

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Platform int32

const (
	Platform_PLATFORM_UNKNOWN       Platform = 0
	Platform_PLATFORM_KUAISHOU      Platform = 1
	Platform_PLATFORM_CHUANSHANJIA  Platform = 2
	Platform_PLATFORM_GUANGDIANTONG Platform = 3
	Platform_PLATFORM_CHUBAO        Platform = 4
	Platform_PLATFORM_QUTOUTIAO     Platform = 5
	Platform_PLATFORM_OTHERS        Platform = 100
)

// Enum value maps for Platform.
var (
	Platform_name = map[int32]string{
		0:   "PLATFORM_UNKNOWN",
		1:   "PLATFORM_KUAISHOU",
		2:   "PLATFORM_CHUANSHANJIA",
		3:   "PLATFORM_GUANGDIANTONG",
		4:   "PLATFORM_CHUBAO",
		5:   "PLATFORM_QUTOUTIAO",
		100: "PLATFORM_OTHERS",
	}
	Platform_value = map[string]int32{
		"PLATFORM_UNKNOWN":       0,
		"PLATFORM_KUAISHOU":      1,
		"PLATFORM_CHUANSHANJIA":  2,
		"PLATFORM_GUANGDIANTONG": 3,
		"PLATFORM_CHUBAO":        4,
		"PLATFORM_QUTOUTIAO":     5,
		"PLATFORM_OTHERS":        100,
	}
)

func (x Platform) Enum() *Platform {
	p := new(Platform)
	*p = x
	return p
}

func (x Platform) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Platform) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_event_hub_proto_enumTypes[0].Descriptor()
}

func (Platform) Type() protoreflect.EnumType {
	return &file_proto_event_hub_proto_enumTypes[0]
}

func (x Platform) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Platform.Descriptor instead.
func (Platform) EnumDescriptor() ([]byte, []int) {
	return file_proto_event_hub_proto_rawDescGZIP(), []int{0}
}

// 29092 common-event-hub-match
type Match struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp       string   `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	PackageName     string   `protobuf:"bytes,2,opt,name=package_name,json=packageName,proto3" json:"package_name,omitempty"`
	Event           string   `protobuf:"bytes,3,opt,name=event,proto3" json:"event,omitempty"`
	UserId          int64    `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Oaid            string   `protobuf:"bytes,5,opt,name=oaid,proto3" json:"oaid,omitempty"`
	Imei            string   `protobuf:"bytes,6,opt,name=imei,proto3" json:"imei,omitempty"`
	Mac             string   `protobuf:"bytes,7,opt,name=mac,proto3" json:"mac,omitempty"`
	AndroidId       string   `protobuf:"bytes,8,opt,name=android_id,json=androidId,proto3" json:"android_id,omitempty"`
	Suuid           string   `protobuf:"bytes,9,opt,name=suuid,proto3" json:"suuid,omitempty"`
	Ip              string   `protobuf:"bytes,10,opt,name=ip,proto3" json:"ip,omitempty"`
	DownloadChannel string   `protobuf:"bytes,11,opt,name=download_channel,json=downloadChannel,proto3" json:"download_channel,omitempty"`
	RegisterDays    int32    `protobuf:"varint,12,opt,name=register_days,json=registerDays,proto3" json:"register_days,omitempty"`
	Topic           string   `protobuf:"bytes,13,opt,name=topic,proto3" json:"topic,omitempty"`
	ClickChannel    string   `protobuf:"bytes,14,opt,name=click_channel,json=clickChannel,proto3" json:"click_channel,omitempty"`
	Cid             string   `protobuf:"bytes,15,opt,name=cid,proto3" json:"cid,omitempty"`
	Aid             string   `protobuf:"bytes,16,opt,name=aid,proto3" json:"aid,omitempty"`
	CampaignId      string   `protobuf:"bytes,17,opt,name=campaign_id,json=campaignId,proto3" json:"campaign_id,omitempty"`
	UnionSite       string   `protobuf:"bytes,18,opt,name=union_site,json=unionSite,proto3" json:"union_site,omitempty"`
	Platform        Platform `protobuf:"varint,19,opt,name=platform,proto3,enum=event_hub.v1.Platform" json:"platform,omitempty"`
	Version         string   `protobuf:"bytes,20,opt,name=version,proto3" json:"version,omitempty"`
	CtaChannel      string   `protobuf:"bytes,21,opt,name=cta_channel,json=ctaChannel,proto3" json:"cta_channel,omitempty"`
}

func (x *Match) Reset() {
	*x = Match{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_hub_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Match) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Match) ProtoMessage() {}

func (x *Match) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_hub_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Match.ProtoReflect.Descriptor instead.
func (*Match) Descriptor() ([]byte, []int) {
	return file_proto_event_hub_proto_rawDescGZIP(), []int{0}
}

func (x *Match) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *Match) GetPackageName() string {
	if x != nil {
		return x.PackageName
	}
	return ""
}

func (x *Match) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *Match) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Match) GetOaid() string {
	if x != nil {
		return x.Oaid
	}
	return ""
}

func (x *Match) GetImei() string {
	if x != nil {
		return x.Imei
	}
	return ""
}

func (x *Match) GetMac() string {
	if x != nil {
		return x.Mac
	}
	return ""
}

func (x *Match) GetAndroidId() string {
	if x != nil {
		return x.AndroidId
	}
	return ""
}

func (x *Match) GetSuuid() string {
	if x != nil {
		return x.Suuid
	}
	return ""
}

func (x *Match) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Match) GetDownloadChannel() string {
	if x != nil {
		return x.DownloadChannel
	}
	return ""
}

func (x *Match) GetRegisterDays() int32 {
	if x != nil {
		return x.RegisterDays
	}
	return 0
}

func (x *Match) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Match) GetClickChannel() string {
	if x != nil {
		return x.ClickChannel
	}
	return ""
}

func (x *Match) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

func (x *Match) GetAid() string {
	if x != nil {
		return x.Aid
	}
	return ""
}

func (x *Match) GetCampaignId() string {
	if x != nil {
		return x.CampaignId
	}
	return ""
}

func (x *Match) GetUnionSite() string {
	if x != nil {
		return x.UnionSite
	}
	return ""
}

func (x *Match) GetPlatform() Platform {
	if x != nil {
		return x.Platform
	}
	return Platform_PLATFORM_UNKNOWN
}

func (x *Match) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Match) GetCtaChannel() string {
	if x != nil {
		return x.CtaChannel
	}
	return ""
}

// 29092 common-event-hub-callback
type Callback struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PackageName     string `protobuf:"bytes,1,opt,name=package_name,json=packageName,proto3" json:"package_name,omitempty"`
	DownloadChannel string `protobuf:"bytes,2,opt,name=download_channel,json=downloadChannel,proto3" json:"download_channel,omitempty"`
	ClickChannel    string `protobuf:"bytes,3,opt,name=click_channel,json=clickChannel,proto3" json:"click_channel,omitempty"`
	CampaignId      string `protobuf:"bytes,4,opt,name=campaign_id,json=campaignId,proto3" json:"campaign_id,omitempty"`
	Aid             string `protobuf:"bytes,5,opt,name=aid,proto3" json:"aid,omitempty"`
	Cid             string `protobuf:"bytes,6,opt,name=cid,proto3" json:"cid,omitempty"`
	Event           string `protobuf:"bytes,7,opt,name=event,proto3" json:"event,omitempty"`
	Timestamp       string `protobuf:"bytes,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Url             string `protobuf:"bytes,9,opt,name=url,proto3" json:"url,omitempty"`
	UserId          string `protobuf:"bytes,10,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Suuid           string `protobuf:"bytes,11,opt,name=suuid,proto3" json:"suuid,omitempty"`
	Version         string `protobuf:"bytes,12,opt,name=version,proto3" json:"version,omitempty"`
	CtaChannel      string `protobuf:"bytes,13,opt,name=cta_channel,json=ctaChannel,proto3" json:"cta_channel,omitempty"`
}

func (x *Callback) Reset() {
	*x = Callback{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_event_hub_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Callback) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Callback) ProtoMessage() {}

func (x *Callback) ProtoReflect() protoreflect.Message {
	mi := &file_proto_event_hub_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Callback.ProtoReflect.Descriptor instead.
func (*Callback) Descriptor() ([]byte, []int) {
	return file_proto_event_hub_proto_rawDescGZIP(), []int{1}
}

func (x *Callback) GetPackageName() string {
	if x != nil {
		return x.PackageName
	}
	return ""
}

func (x *Callback) GetDownloadChannel() string {
	if x != nil {
		return x.DownloadChannel
	}
	return ""
}

func (x *Callback) GetClickChannel() string {
	if x != nil {
		return x.ClickChannel
	}
	return ""
}

func (x *Callback) GetCampaignId() string {
	if x != nil {
		return x.CampaignId
	}
	return ""
}

func (x *Callback) GetAid() string {
	if x != nil {
		return x.Aid
	}
	return ""
}

func (x *Callback) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

func (x *Callback) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *Callback) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *Callback) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Callback) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Callback) GetSuuid() string {
	if x != nil {
		return x.Suuid
	}
	return ""
}

func (x *Callback) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Callback) GetCtaChannel() string {
	if x != nil {
		return x.CtaChannel
	}
	return ""
}

var File_proto_event_hub_proto protoreflect.FileDescriptor

var file_proto_event_hub_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x75,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x68,
	0x75, 0x62, 0x2e, 0x76, 0x31, 0x22, 0xd4, 0x04, 0x0a, 0x05, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x21, 0x0a,
	0x0c, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6f, 0x61, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6f,
	0x61, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6d, 0x65, 0x69, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x69, 0x6d, 0x65, 0x69, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x63, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x61, 0x63, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x6e, 0x64,
	0x72, 0x6f, 0x69, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61,
	0x6e, 0x64, 0x72, 0x6f, 0x69, 0x64, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x75, 0x75, 0x69, 0x64, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x29,
	0x0a, 0x10, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f,
	0x61, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x79, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x44, 0x61, 0x79, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x63, 0x6b, 0x5f, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69,
	0x63, 0x6b, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x69, 0x64,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x61,
	0x69, 0x64, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x69, 0x64, 0x12, 0x1f, 0x0a,
	0x0b, 0x63, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x11, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x49, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x75, 0x6e, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x69, 0x74, 0x65, 0x18, 0x12, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x75, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x69, 0x74, 0x65, 0x12, 0x32, 0x0a,
	0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x14, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x74, 0x61, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x63, 0x74, 0x61, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x22, 0xf2, 0x02, 0x0a,
	0x08, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x61, 0x63,
	0x6b, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x10,
	0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x63, 0x6b,
	0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x63, 0x6c, 0x69, 0x63, 0x6b, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1f, 0x0a, 0x0b,
	0x63, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x63, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x49, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x61, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x69, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x63, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x73, 0x75, 0x75, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x73, 0x75, 0x75, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x74, 0x61, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x74, 0x61, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x2a, 0xb0, 0x01, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x14,
	0x0a, 0x10, 0x50, 0x4c, 0x41, 0x54, 0x46, 0x4f, 0x52, 0x4d, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f,
	0x57, 0x4e, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x4c, 0x41, 0x54, 0x46, 0x4f, 0x52, 0x4d,
	0x5f, 0x4b, 0x55, 0x41, 0x49, 0x53, 0x48, 0x4f, 0x55, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x50,
	0x4c, 0x41, 0x54, 0x46, 0x4f, 0x52, 0x4d, 0x5f, 0x43, 0x48, 0x55, 0x41, 0x4e, 0x53, 0x48, 0x41,
	0x4e, 0x4a, 0x49, 0x41, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x50, 0x4c, 0x41, 0x54, 0x46, 0x4f,
	0x52, 0x4d, 0x5f, 0x47, 0x55, 0x41, 0x4e, 0x47, 0x44, 0x49, 0x41, 0x4e, 0x54, 0x4f, 0x4e, 0x47,
	0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x4c, 0x41, 0x54, 0x46, 0x4f, 0x52, 0x4d, 0x5f, 0x43,
	0x48, 0x55, 0x42, 0x41, 0x4f, 0x10, 0x04, 0x12, 0x16, 0x0a, 0x12, 0x50, 0x4c, 0x41, 0x54, 0x46,
	0x4f, 0x52, 0x4d, 0x5f, 0x51, 0x55, 0x54, 0x4f, 0x55, 0x54, 0x49, 0x41, 0x4f, 0x10, 0x05, 0x12,
	0x13, 0x0a, 0x0f, 0x50, 0x4c, 0x41, 0x54, 0x46, 0x4f, 0x52, 0x4d, 0x5f, 0x4f, 0x54, 0x48, 0x45,
	0x52, 0x53, 0x10, 0x64, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x6c, 0x61, 0x62, 0x2e, 0x74, 0x61, 0x67,
	0x74, 0x69, 0x63, 0x2e, 0x63, 0x6e, 0x2f, 0x61, 0x64, 0x5f, 0x67, 0x61, 0x69, 0x6e, 0x73, 0x2f,
	0x6b, 0x69, 0x74, 0x74, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x6b, 0x69, 0x74, 0x74,
	0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_event_hub_proto_rawDescOnce sync.Once
	file_proto_event_hub_proto_rawDescData = file_proto_event_hub_proto_rawDesc
)

func file_proto_event_hub_proto_rawDescGZIP() []byte {
	file_proto_event_hub_proto_rawDescOnce.Do(func() {
		file_proto_event_hub_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_event_hub_proto_rawDescData)
	})
	return file_proto_event_hub_proto_rawDescData
}

var file_proto_event_hub_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_event_hub_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_event_hub_proto_goTypes = []interface{}{
	(Platform)(0),    // 0: event_hub.v1.Platform
	(*Match)(nil),    // 1: event_hub.v1.Match
	(*Callback)(nil), // 2: event_hub.v1.Callback
}
var file_proto_event_hub_proto_depIdxs = []int32{
	0, // 0: event_hub.v1.Match.platform:type_name -> event_hub.v1.Platform
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_event_hub_proto_init() }
func file_proto_event_hub_proto_init() {
	if File_proto_event_hub_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_event_hub_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Match); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_event_hub_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Callback); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_event_hub_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_event_hub_proto_goTypes,
		DependencyIndexes: file_proto_event_hub_proto_depIdxs,
		EnumInfos:         file_proto_event_hub_proto_enumTypes,
		MessageInfos:      file_proto_event_hub_proto_msgTypes,
	}.Build()
	File_proto_event_hub_proto = out.File
	file_proto_event_hub_proto_rawDesc = nil
	file_proto_event_hub_proto_goTypes = nil
	file_proto_event_hub_proto_depIdxs = nil
}
