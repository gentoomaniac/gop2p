// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: protofiles/messages.proto

package main

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    int64  `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Payload string `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_protofiles_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Message) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

// Peer management
type GetPeers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit int64 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetPeers) Reset() {
	*x = GetPeers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPeers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPeers) ProtoMessage() {}

func (x *GetPeers) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPeers.ProtoReflect.Descriptor instead.
func (*GetPeers) Descriptor() ([]byte, []int) {
	return file_protofiles_messages_proto_rawDescGZIP(), []int{1}
}

func (x *GetPeers) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type SendPeers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumRecords int64 `protobuf:"varint,1,opt,name=numRecords,proto3" json:"numRecords,omitempty"`
}

func (x *SendPeers) Reset() {
	*x = SendPeers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendPeers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendPeers) ProtoMessage() {}

func (x *SendPeers) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendPeers.ProtoReflect.Descriptor instead.
func (*SendPeers) Descriptor() ([]byte, []int) {
	return file_protofiles_messages_proto_rawDescGZIP(), []int{2}
}

func (x *SendPeers) GetNumRecords() int64 {
	if x != nil {
		return x.NumRecords
	}
	return 0
}

type PbPeer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address  string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Port     int64  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Protocol string `protobuf:"bytes,3,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *PbPeer) Reset() {
	*x = PbPeer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbPeer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbPeer) ProtoMessage() {}

func (x *PbPeer) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbPeer.ProtoReflect.Descriptor instead.
func (*PbPeer) Descriptor() ([]byte, []int) {
	return file_protofiles_messages_proto_rawDescGZIP(), []int{3}
}

func (x *PbPeer) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *PbPeer) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *PbPeer) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

var File_protofiles_messages_proto protoreflect.FileDescriptor

var file_protofiles_messages_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69,
	0x6e, 0x22, 0x37, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x20, 0x0a, 0x08, 0x47, 0x65,
	0x74, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x2b, 0x0a, 0x09,
	0x53, 0x65, 0x6e, 0x64, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x6e, 0x75, 0x6d,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6e,
	0x75, 0x6d, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x52, 0x0a, 0x06, 0x50, 0x62, 0x50,
	0x65, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x42, 0x1f, 0x5a,
	0x1d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x6e, 0x74,
	0x6f, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x61, 0x63, 0x2f, 0x67, 0x6f, 0x70, 0x32, 0x70, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protofiles_messages_proto_rawDescOnce sync.Once
	file_protofiles_messages_proto_rawDescData = file_protofiles_messages_proto_rawDesc
)

func file_protofiles_messages_proto_rawDescGZIP() []byte {
	file_protofiles_messages_proto_rawDescOnce.Do(func() {
		file_protofiles_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_protofiles_messages_proto_rawDescData)
	})
	return file_protofiles_messages_proto_rawDescData
}

var file_protofiles_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protofiles_messages_proto_goTypes = []interface{}{
	(*Message)(nil),   // 0: main.Message
	(*GetPeers)(nil),  // 1: main.GetPeers
	(*SendPeers)(nil), // 2: main.SendPeers
	(*PbPeer)(nil),    // 3: main.PbPeer
}
var file_protofiles_messages_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protofiles_messages_proto_init() }
func file_protofiles_messages_proto_init() {
	if File_protofiles_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protofiles_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_protofiles_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPeers); i {
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
		file_protofiles_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendPeers); i {
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
		file_protofiles_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbPeer); i {
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
			RawDescriptor: file_protofiles_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protofiles_messages_proto_goTypes,
		DependencyIndexes: file_protofiles_messages_proto_depIdxs,
		MessageInfos:      file_protofiles_messages_proto_msgTypes,
	}.Build()
	File_protofiles_messages_proto = out.File
	file_protofiles_messages_proto_rawDesc = nil
	file_protofiles_messages_proto_goTypes = nil
	file_protofiles_messages_proto_depIdxs = nil
}