// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: c2s.proto

package protocol

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

type ID int32

const (
	ID_C2S_LOGIN_REQ   ID = 1001 // C2SLoginReq
	ID_S2C_LOGIN_RESP  ID = 1002 // S2CLoginResp
	ID_C2S_LOGOUT_REQ  ID = 1003 // C2SLogoutReq
	ID_S2C_LOGOUT_RESP ID = 1004 // S2CLogoutResp
)

// Enum value maps for ID.
var (
	ID_name = map[int32]string{
		1001: "C2S_LOGIN_REQ",
		1002: "S2C_LOGIN_RESP",
		1003: "C2S_LOGOUT_REQ",
		1004: "S2C_LOGOUT_RESP",
	}
	ID_value = map[string]int32{
		"C2S_LOGIN_REQ":   1001,
		"S2C_LOGIN_RESP":  1002,
		"C2S_LOGOUT_REQ":  1003,
		"S2C_LOGOUT_RESP": 1004,
	}
)

func (x ID) Enum() *ID {
	p := new(ID)
	*p = x
	return p
}

func (x ID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ID) Descriptor() protoreflect.EnumDescriptor {
	return file_c2s_proto_enumTypes[0].Descriptor()
}

func (ID) Type() protoreflect.EnumType {
	return &file_c2s_proto_enumTypes[0]
}

func (x ID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *ID) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = ID(num)
	return nil
}

// Deprecated: Use ID.Descriptor instead.
func (ID) EnumDescriptor() ([]byte, []int) {
	return file_c2s_proto_rawDescGZIP(), []int{0}
}

type C2SLoginReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Account *string `protobuf:"bytes,1,opt,name=account" json:"account,omitempty"`
	Age     *int32  `protobuf:"varint,2,opt,name=age" json:"age,omitempty"`
	Friends []int32 `protobuf:"varint,3,rep,name=friends" json:"friends,omitempty"`
}

func (x *C2SLoginReq) Reset() {
	*x = C2SLoginReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_c2s_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SLoginReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SLoginReq) ProtoMessage() {}

func (x *C2SLoginReq) ProtoReflect() protoreflect.Message {
	mi := &file_c2s_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SLoginReq.ProtoReflect.Descriptor instead.
func (*C2SLoginReq) Descriptor() ([]byte, []int) {
	return file_c2s_proto_rawDescGZIP(), []int{0}
}

func (x *C2SLoginReq) GetAccount() string {
	if x != nil && x.Account != nil {
		return *x.Account
	}
	return ""
}

func (x *C2SLoginReq) GetAge() int32 {
	if x != nil && x.Age != nil {
		return *x.Age
	}
	return 0
}

func (x *C2SLoginReq) GetFriends() []int32 {
	if x != nil {
		return x.Friends
	}
	return nil
}

type S2CLoginResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ret *int32 `protobuf:"varint,1,opt,name=ret" json:"ret,omitempty"`
}

func (x *S2CLoginResp) Reset() {
	*x = S2CLoginResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_c2s_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CLoginResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CLoginResp) ProtoMessage() {}

func (x *S2CLoginResp) ProtoReflect() protoreflect.Message {
	mi := &file_c2s_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CLoginResp.ProtoReflect.Descriptor instead.
func (*S2CLoginResp) Descriptor() ([]byte, []int) {
	return file_c2s_proto_rawDescGZIP(), []int{1}
}

func (x *S2CLoginResp) GetRet() int32 {
	if x != nil && x.Ret != nil {
		return *x.Ret
	}
	return 0
}

type C2SLogoutReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *C2SLogoutReq) Reset() {
	*x = C2SLogoutReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_c2s_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2SLogoutReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2SLogoutReq) ProtoMessage() {}

func (x *C2SLogoutReq) ProtoReflect() protoreflect.Message {
	mi := &file_c2s_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2SLogoutReq.ProtoReflect.Descriptor instead.
func (*C2SLogoutReq) Descriptor() ([]byte, []int) {
	return file_c2s_proto_rawDescGZIP(), []int{2}
}

type S2CLogoutResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ret *int32 `protobuf:"varint,1,opt,name=ret" json:"ret,omitempty"`
}

func (x *S2CLogoutResp) Reset() {
	*x = S2CLogoutResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_c2s_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2CLogoutResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2CLogoutResp) ProtoMessage() {}

func (x *S2CLogoutResp) ProtoReflect() protoreflect.Message {
	mi := &file_c2s_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2CLogoutResp.ProtoReflect.Descriptor instead.
func (*S2CLogoutResp) Descriptor() ([]byte, []int) {
	return file_c2s_proto_rawDescGZIP(), []int{3}
}

func (x *S2CLogoutResp) GetRet() int32 {
	if x != nil && x.Ret != nil {
		return *x.Ret
	}
	return 0
}

var File_c2s_proto protoreflect.FileDescriptor

var file_c2s_proto_rawDesc = []byte{
	0x0a, 0x09, 0x63, 0x32, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x53, 0x0a, 0x0b, 0x43, 0x32, 0x53, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x61, 0x67, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x05, 0x52, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x22, 0x20, 0x0a, 0x0c, 0x53, 0x32,
	0x43, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x72, 0x65, 0x74, 0x22, 0x0e, 0x0a, 0x0c,
	0x43, 0x32, 0x53, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x22, 0x21, 0x0a, 0x0d,
	0x53, 0x32, 0x43, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a,
	0x03, 0x72, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x72, 0x65, 0x74, 0x2a,
	0x58, 0x0a, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x0d, 0x43, 0x32, 0x53, 0x5f, 0x4c, 0x4f, 0x47,
	0x49, 0x4e, 0x5f, 0x52, 0x45, 0x51, 0x10, 0xe9, 0x07, 0x12, 0x13, 0x0a, 0x0e, 0x53, 0x32, 0x43,
	0x5f, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x10, 0xea, 0x07, 0x12, 0x13,
	0x0a, 0x0e, 0x43, 0x32, 0x53, 0x5f, 0x4c, 0x4f, 0x47, 0x4f, 0x55, 0x54, 0x5f, 0x52, 0x45, 0x51,
	0x10, 0xeb, 0x07, 0x12, 0x14, 0x0a, 0x0f, 0x53, 0x32, 0x43, 0x5f, 0x4c, 0x4f, 0x47, 0x4f, 0x55,
	0x54, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x10, 0xec, 0x07, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x3b,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
}

var (
	file_c2s_proto_rawDescOnce sync.Once
	file_c2s_proto_rawDescData = file_c2s_proto_rawDesc
)

func file_c2s_proto_rawDescGZIP() []byte {
	file_c2s_proto_rawDescOnce.Do(func() {
		file_c2s_proto_rawDescData = protoimpl.X.CompressGZIP(file_c2s_proto_rawDescData)
	})
	return file_c2s_proto_rawDescData
}

var file_c2s_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_c2s_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_c2s_proto_goTypes = []interface{}{
	(ID)(0),               // 0: protocol.ID
	(*C2SLoginReq)(nil),   // 1: protocol.C2SLoginReq
	(*S2CLoginResp)(nil),  // 2: protocol.S2CLoginResp
	(*C2SLogoutReq)(nil),  // 3: protocol.C2SLogoutReq
	(*S2CLogoutResp)(nil), // 4: protocol.S2CLogoutResp
}
var file_c2s_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_c2s_proto_init() }
func file_c2s_proto_init() {
	if File_c2s_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_c2s_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SLoginReq); i {
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
		file_c2s_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CLoginResp); i {
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
		file_c2s_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2SLogoutReq); i {
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
		file_c2s_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2CLogoutResp); i {
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
			RawDescriptor: file_c2s_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_c2s_proto_goTypes,
		DependencyIndexes: file_c2s_proto_depIdxs,
		EnumInfos:         file_c2s_proto_enumTypes,
		MessageInfos:      file_c2s_proto_msgTypes,
	}.Build()
	File_c2s_proto = out.File
	file_c2s_proto_rawDesc = nil
	file_c2s_proto_goTypes = nil
	file_c2s_proto_depIdxs = nil
}