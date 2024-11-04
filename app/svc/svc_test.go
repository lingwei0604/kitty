package svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	pb "git.yingzhongshare.com/mkt/kitty/proto"
	"github.com/gogo/protobuf/jsonpb"
)

type GenericReply struct {
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	// deprecated
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Msg                  string   `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func TestMarshal(t *testing.T) {
	foo := pb.GenericReply{
		Code:    0,
		Message: "",
		Msg:     "",
	}
	marshaller := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "",
		OrigName:     true,
		AnyResolver:  nil,
	}

	var str bytes.Buffer
	marshaller.Marshal(&str, &foo)
	fmt.Println(str.String())

	s, _ := json.Marshal(&foo)
	fmt.Println(string(s))
}
