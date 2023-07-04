// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: storage_message.proto

package pb

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

type Storage_Drive int32

const (
	Storage_UNKNOWN Storage_Drive = 0
	Storage_HDD     Storage_Drive = 1
	Storage_SSD     Storage_Drive = 2
)

// Enum value maps for Storage_Drive.
var (
	Storage_Drive_name = map[int32]string{
		0: "UNKNOWN",
		1: "HDD",
		2: "SSD",
	}
	Storage_Drive_value = map[string]int32{
		"UNKNOWN": 0,
		"HDD":     1,
		"SSD":     2,
	}
)

func (x Storage_Drive) Enum() *Storage_Drive {
	p := new(Storage_Drive)
	*p = x
	return p
}

func (x Storage_Drive) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Storage_Drive) Descriptor() protoreflect.EnumDescriptor {
	return file_storage_message_proto_enumTypes[0].Descriptor()
}

func (Storage_Drive) Type() protoreflect.EnumType {
	return &file_storage_message_proto_enumTypes[0]
}

func (x Storage_Drive) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Storage_Drive.Descriptor instead.
func (Storage_Drive) EnumDescriptor() ([]byte, []int) {
	return file_storage_message_proto_rawDescGZIP(), []int{0, 0}
}

type Storage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Drive  Storage_Drive `protobuf:"varint,1,opt,name=drive,proto3,enum=pb.Storage_Drive" json:"drive,omitempty"`
	Memory *Memory       `protobuf:"bytes,2,opt,name=memory,proto3" json:"memory,omitempty"`
}

func (x *Storage) Reset() {
	*x = Storage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Storage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Storage) ProtoMessage() {}

func (x *Storage) ProtoReflect() protoreflect.Message {
	mi := &file_storage_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Storage.ProtoReflect.Descriptor instead.
func (*Storage) Descriptor() ([]byte, []int) {
	return file_storage_message_proto_rawDescGZIP(), []int{0}
}

func (x *Storage) GetDrive() Storage_Drive {
	if x != nil {
		return x.Drive
	}
	return Storage_UNKNOWN
}

func (x *Storage) GetMemory() *Memory {
	if x != nil {
		return x.Memory
	}
	return nil
}

var File_storage_message_proto protoreflect.FileDescriptor

var file_storage_message_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x14, 0x6d, 0x65, 0x6d,
	0x6f, 0x72, 0x79, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x7e, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x27, 0x0a, 0x05,
	0x64, 0x72, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x62,
	0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x72, 0x69, 0x76, 0x65, 0x52, 0x05,
	0x64, 0x72, 0x69, 0x76, 0x65, 0x12, 0x22, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x6d, 0x6f, 0x72,
	0x79, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x22, 0x26, 0x0a, 0x05, 0x44, 0x72, 0x69,
	0x76, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x07, 0x0a, 0x03, 0x48, 0x44, 0x44, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x53, 0x44, 0x10,
	0x02, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storage_message_proto_rawDescOnce sync.Once
	file_storage_message_proto_rawDescData = file_storage_message_proto_rawDesc
)

func file_storage_message_proto_rawDescGZIP() []byte {
	file_storage_message_proto_rawDescOnce.Do(func() {
		file_storage_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_storage_message_proto_rawDescData)
	})
	return file_storage_message_proto_rawDescData
}

var file_storage_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_storage_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_storage_message_proto_goTypes = []interface{}{
	(Storage_Drive)(0), // 0: pb.Storage.Drive
	(*Storage)(nil),    // 1: pb.Storage
	(*Memory)(nil),     // 2: pb.Memory
}
var file_storage_message_proto_depIdxs = []int32{
	0, // 0: pb.Storage.drive:type_name -> pb.Storage.Drive
	2, // 1: pb.Storage.memory:type_name -> pb.Memory
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_storage_message_proto_init() }
func file_storage_message_proto_init() {
	if File_storage_message_proto != nil {
		return
	}
	file_memory_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_storage_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Storage); i {
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
			RawDescriptor: file_storage_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_storage_message_proto_goTypes,
		DependencyIndexes: file_storage_message_proto_depIdxs,
		EnumInfos:         file_storage_message_proto_enumTypes,
		MessageInfos:      file_storage_message_proto_msgTypes,
	}.Build()
	File_storage_message_proto = out.File
	file_storage_message_proto_rawDesc = nil
	file_storage_message_proto_goTypes = nil
	file_storage_message_proto_depIdxs = nil
}
