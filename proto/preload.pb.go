// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: preload.proto

package kitty

import (
	context "context"
	fmt "fmt"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PreloadReq struct {
	PreloadHostList      []string `protobuf:"bytes,1,rep,name=preloadHostList,proto3" json:"preloadHostList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PreloadReq) Reset()         { *m = PreloadReq{} }
func (m *PreloadReq) String() string { return proto.CompactTextString(m) }
func (*PreloadReq) ProtoMessage()    {}
func (*PreloadReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf9bffd16cc7e6f9, []int{0}
}
func (m *PreloadReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PreloadReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PreloadReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PreloadReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PreloadReq.Merge(m, src)
}
func (m *PreloadReq) XXX_Size() int {
	return m.Size()
}
func (m *PreloadReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PreloadReq.DiscardUnknown(m)
}

var xxx_messageInfo_PreloadReq proto.InternalMessageInfo

func (m *PreloadReq) GetPreloadHostList() []string {
	if m != nil {
		return m.PreloadHostList
	}
	return nil
}

type PreloadResp struct {
	Code                 int32          `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string         `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data                 []*PreloadInfo `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PreloadResp) Reset()         { *m = PreloadResp{} }
func (m *PreloadResp) String() string { return proto.CompactTextString(m) }
func (*PreloadResp) ProtoMessage()    {}
func (*PreloadResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf9bffd16cc7e6f9, []int{1}
}
func (m *PreloadResp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PreloadResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PreloadResp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PreloadResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PreloadResp.Merge(m, src)
}
func (m *PreloadResp) XXX_Size() int {
	return m.Size()
}
func (m *PreloadResp) XXX_DiscardUnknown() {
	xxx_messageInfo_PreloadResp.DiscardUnknown(m)
}

var xxx_messageInfo_PreloadResp proto.InternalMessageInfo

func (m *PreloadResp) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *PreloadResp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *PreloadResp) GetData() []*PreloadInfo {
	if m != nil {
		return m.Data
	}
	return nil
}

type PreloadInfo struct {
	Gzurl                string   `protobuf:"bytes,1,opt,name=gzurl,proto3" json:"gzurl,omitempty"`
	Md5                  string   `protobuf:"bytes,2,opt,name=md5,proto3" json:"md5,omitempty"`
	Weburl               string   `protobuf:"bytes,3,opt,name=weburl,proto3" json:"weburl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PreloadInfo) Reset()         { *m = PreloadInfo{} }
func (m *PreloadInfo) String() string { return proto.CompactTextString(m) }
func (*PreloadInfo) ProtoMessage()    {}
func (*PreloadInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf9bffd16cc7e6f9, []int{2}
}
func (m *PreloadInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PreloadInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PreloadInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PreloadInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PreloadInfo.Merge(m, src)
}
func (m *PreloadInfo) XXX_Size() int {
	return m.Size()
}
func (m *PreloadInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PreloadInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PreloadInfo proto.InternalMessageInfo

func (m *PreloadInfo) GetGzurl() string {
	if m != nil {
		return m.Gzurl
	}
	return ""
}

func (m *PreloadInfo) GetMd5() string {
	if m != nil {
		return m.Md5
	}
	return ""
}

func (m *PreloadInfo) GetWeburl() string {
	if m != nil {
		return m.Weburl
	}
	return ""
}

func init() {
	proto.RegisterType((*PreloadReq)(nil), "preload.v1.PreloadReq")
	proto.RegisterType((*PreloadResp)(nil), "preload.v1.PreloadResp")
	proto.RegisterType((*PreloadInfo)(nil), "preload.v1.PreloadInfo")
}

func init() { proto.RegisterFile("preload.proto", fileDescriptor_cf9bffd16cc7e6f9) }

var fileDescriptor_cf9bffd16cc7e6f9 = []byte{
	// 551 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xcf, 0x6b, 0x13, 0x41,
	0x18, 0x75, 0xf3, 0xcb, 0xee, 0x84, 0x62, 0x1c, 0x6b, 0x13, 0x42, 0x09, 0x61, 0xbd, 0x44, 0x4a,
	0x76, 0x9a, 0x48, 0x2f, 0xf1, 0x62, 0x73, 0x10, 0x05, 0x15, 0x59, 0xbc, 0xe8, 0x45, 0x27, 0xbb,
	0xdb, 0x75, 0x74, 0x33, 0x33, 0xdd, 0x1d, 0xd3, 0xda, 0xa3, 0x14, 0x85, 0xe2, 0x2f, 0x22, 0x06,
	0x11, 0x05, 0xab, 0x9e, 0x04, 0xa1, 0x17, 0x21, 0x08, 0xbd, 0x7b, 0x14, 0xfc, 0x07, 0x24, 0x66,
	0x97, 0xfa, 0x17, 0x08, 0x3d, 0xc9, 0x4e, 0xd6, 0x58, 0xda, 0x9c, 0xf6, 0xbd, 0xc7, 0xfb, 0xde,
	0xfb, 0x76, 0xbf, 0x05, 0xd3, 0xdc, 0xb3, 0x5d, 0x86, 0x2d, 0x9d, 0x7b, 0x4c, 0x30, 0x08, 0xfe,
	0xd1, 0x4e, 0xad, 0x38, 0xe7, 0x30, 0xe6, 0xb8, 0x36, 0xc2, 0x9c, 0x20, 0x4c, 0x29, 0x13, 0x58,
	0x10, 0x46, 0xfd, 0x91, 0xb3, 0x38, 0x2f, 0x1f, 0x66, 0xd5, 0xb1, 0x69, 0xb5, 0x83, 0x5d, 0x62,
	0x61, 0x61, 0xa3, 0x43, 0x20, 0x36, 0xeb, 0xfb, 0xcc, 0x8c, 0xdb, 0x14, 0x73, 0xd2, 0xa9, 0x23,
	0xc6, 0x65, 0xe0, 0x84, 0xf0, 0x19, 0x87, 0x39, 0x4c, 0x42, 0x14, 0xa1, 0x91, 0xaa, 0x5d, 0x03,
	0xe0, 0xea, 0x68, 0x3d, 0xc3, 0x5e, 0x81, 0xe7, 0xc1, 0xb1, 0x78, 0xd9, 0x0b, 0xcc, 0x17, 0x97,
	0x88, 0x2f, 0x0a, 0x4a, 0x39, 0x59, 0x51, 0x9b, 0x73, 0x7b, 0x4d, 0xb5, 0xab, 0x64, 0xb4, 0x94,
	0x97, 0xc8, 0x25, 0xf6, 0x9a, 0xd3, 0x5d, 0x05, 0x68, 0x53, 0x5e, 0xa6, 0x91, 0xba, 0x2d, 0x04,
	0x37, 0x0e, 0x0e, 0x69, 0xb7, 0x40, 0x76, 0x9c, 0xea, 0x73, 0x08, 0x41, 0xca, 0x64, 0x96, 0x5d,
	0x50, 0xca, 0x4a, 0x25, 0x6d, 0x48, 0x0c, 0x73, 0x20, 0xd9, 0xf6, 0x9d, 0x42, 0xa2, 0xac, 0x54,
	0x54, 0x23, 0x82, 0x70, 0x1e, 0xa4, 0x2c, 0x2c, 0x70, 0x21, 0x59, 0x4e, 0x56, 0xb2, 0xf5, 0xbc,
	0xfe, 0xff, 0xb3, 0xe9, 0x71, 0xd8, 0x45, 0xba, 0xcc, 0x0c, 0x69, 0xd2, 0x2e, 0x8f, 0x1b, 0x22,
	0x11, 0xce, 0x80, 0xb4, 0xb3, 0x7e, 0xcf, 0x73, 0x65, 0x85, 0x6a, 0x8c, 0x88, 0xec, 0xb0, 0x16,
	0xc7, 0x1d, 0xd6, 0x22, 0x9c, 0x05, 0x99, 0x55, 0xbb, 0x15, 0x19, 0x93, 0x52, 0x8c, 0x59, 0xfd,
	0x3a, 0x38, 0x1a, 0xc7, 0xc1, 0x2b, 0x60, 0x2a, 0x7a, 0x07, 0x19, 0x3b, 0x3b, 0x61, 0x09, 0xc3,
	0x5e, 0x29, 0xe6, 0x27, 0xea, 0x3e, 0xd7, 0x8e, 0x3f, 0xf8, 0x31, 0x7c, 0x91, 0xc8, 0x42, 0x15,
	0x11, 0xba, 0xcc, 0xa2, 0x9c, 0xe6, 0x1f, 0xa5, 0xbb, 0xd4, 0x57, 0x60, 0x13, 0xe4, 0x77, 0xbf,
	0x6c, 0x84, 0x5f, 0x9f, 0xff, 0x7e, 0xbf, 0x11, 0x6c, 0xbe, 0xdb, 0xdd, 0xd9, 0x09, 0xdf, 0x3c,
	0x09, 0x36, 0xb7, 0x3a, 0x35, 0xad, 0x0c, 0xd4, 0xf0, 0xed, 0xd3, 0xe0, 0xe3, 0xa7, 0xb0, 0xb7,
	0x5d, 0x3c, 0xe1, 0x12, 0x81, 0xc9, 0x1a, 0xa1, 0xe7, 0x2c, 0x46, 0xed, 0x55, 0x5f, 0x37, 0x59,
	0xbb, 0x9e, 0xae, 0xe9, 0x0b, 0xfa, 0x82, 0x96, 0x45, 0x71, 0x2d, 0xea, 0xd4, 0xea, 0x39, 0xcc,
	0xb9, 0x4b, 0x4c, 0x79, 0x63, 0x74, 0xc7, 0x67, 0xb4, 0x71, 0x48, 0xf1, 0x96, 0xc0, 0xc9, 0xe0,
	0xe5, 0xc3, 0x61, 0xff, 0x75, 0xf0, 0xb8, 0x37, 0xec, 0x7f, 0x0e, 0x9e, 0x6d, 0x87, 0xbd, 0x47,
	0xe1, 0xd6, 0x07, 0x58, 0x89, 0xee, 0xd6, 0x40, 0xa8, 0xcd, 0xa8, 0x2d, 0xc8, 0xba, 0x1c, 0xa9,
	0x5a, 0xcc, 0xf4, 0xf5, 0x35, 0x47, 0x17, 0xd8, 0x11, 0xc4, 0xd4, 0x4d, 0x8a, 0x4e, 0xa1, 0x6f,
	0x83, 0x92, 0xf2, 0x7d, 0x50, 0x52, 0x7e, 0x0e, 0x4a, 0xca, 0xab, 0x5f, 0xa5, 0x23, 0x37, 0x4e,
	0x3b, 0x2e, 0x6e, 0xed, 0xb3, 0x60, 0xeb, 0xa6, 0x83, 0x09, 0xf5, 0xd1, 0x5d, 0x22, 0xc4, 0x7d,
	0x24, 0xff, 0xaa, 0xb3, 0x12, 0xb7, 0x32, 0x92, 0x9c, 0xf9, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x3d,
	0x7d, 0x29, 0x89, 0x10, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PreloadClient is the client API for Preload service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PreloadClient interface {
	// 获取 preload 预加载资源地址
	ListInfo(ctx context.Context, in *PreloadReq, opts ...grpc.CallOption) (*PreloadResp, error)
}

type preloadClient struct {
	cc *grpc.ClientConn
}

func NewPreloadClient(cc *grpc.ClientConn) PreloadClient {
	return &preloadClient{cc}
}

func (c *preloadClient) ListInfo(ctx context.Context, in *PreloadReq, opts ...grpc.CallOption) (*PreloadResp, error) {
	out := new(PreloadResp)
	err := c.cc.Invoke(ctx, "/preload.v1.Preload/ListInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PreloadServer is the server API for Preload service.
type PreloadServer interface {
	// 获取 preload 预加载资源地址
	ListInfo(context.Context, *PreloadReq) (*PreloadResp, error)
}

// UnimplementedPreloadServer can be embedded to have forward compatible implementations.
type UnimplementedPreloadServer struct {
}

func (*UnimplementedPreloadServer) ListInfo(ctx context.Context, req *PreloadReq) (*PreloadResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInfo not implemented")
}

func RegisterPreloadServer(s *grpc.Server, srv PreloadServer) {
	s.RegisterService(&_Preload_serviceDesc, srv)
}

func _Preload_ListInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PreloadReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PreloadServer).ListInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/preload.v1.Preload/ListInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PreloadServer).ListInfo(ctx, req.(*PreloadReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Preload_serviceDesc = grpc.ServiceDesc{
	ServiceName: "preload.v1.Preload",
	HandlerType: (*PreloadServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListInfo",
			Handler:    _Preload_ListInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "preload.proto",
}

func (m *PreloadReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PreloadReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PreloadReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.PreloadHostList) > 0 {
		for iNdEx := len(m.PreloadHostList) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.PreloadHostList[iNdEx])
			copy(dAtA[i:], m.PreloadHostList[iNdEx])
			i = encodeVarintPreload(dAtA, i, uint64(len(m.PreloadHostList[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *PreloadResp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PreloadResp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PreloadResp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Data) > 0 {
		for iNdEx := len(m.Data) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Data[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPreload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintPreload(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x12
	}
	if m.Code != 0 {
		i = encodeVarintPreload(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PreloadInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PreloadInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PreloadInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Weburl) > 0 {
		i -= len(m.Weburl)
		copy(dAtA[i:], m.Weburl)
		i = encodeVarintPreload(dAtA, i, uint64(len(m.Weburl)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Md5) > 0 {
		i -= len(m.Md5)
		copy(dAtA[i:], m.Md5)
		i = encodeVarintPreload(dAtA, i, uint64(len(m.Md5)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Gzurl) > 0 {
		i -= len(m.Gzurl)
		copy(dAtA[i:], m.Gzurl)
		i = encodeVarintPreload(dAtA, i, uint64(len(m.Gzurl)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPreload(dAtA []byte, offset int, v uint64) int {
	offset -= sovPreload(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PreloadReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.PreloadHostList) > 0 {
		for _, s := range m.PreloadHostList {
			l = len(s)
			n += 1 + l + sovPreload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PreloadResp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovPreload(uint64(m.Code))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovPreload(uint64(l))
	}
	if len(m.Data) > 0 {
		for _, e := range m.Data {
			l = e.Size()
			n += 1 + l + sovPreload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PreloadInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Gzurl)
	if l > 0 {
		n += 1 + l + sovPreload(uint64(l))
	}
	l = len(m.Md5)
	if l > 0 {
		n += 1 + l + sovPreload(uint64(l))
	}
	l = len(m.Weburl)
	if l > 0 {
		n += 1 + l + sovPreload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovPreload(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPreload(x uint64) (n int) {
	return sovPreload(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PreloadReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPreload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PreloadReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PreloadReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreloadHostList", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPreload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPreload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PreloadHostList = append(m.PreloadHostList, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPreload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPreload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPreload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PreloadResp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPreload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PreloadResp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PreloadResp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPreload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPreload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPreload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPreload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data, &PreloadInfo{})
			if err := m.Data[len(m.Data)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPreload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPreload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPreload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PreloadInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPreload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PreloadInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PreloadInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gzurl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPreload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPreload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Gzurl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Md5", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPreload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPreload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Md5 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weburl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPreload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPreload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Weburl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPreload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPreload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPreload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPreload(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPreload
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPreload
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPreload
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPreload
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPreload
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPreload        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPreload          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPreload = fmt.Errorf("proto: unexpected end of group")
)