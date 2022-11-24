// Sorapointa - A server software re-implementation for a certain anime game, and avoid sorapointa.
// Copyright (C) 2022  Sorapointa Team
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0
// source: HomeBlueprintBriefMuipData.proto

package proto

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

type HomeBlueprintBriefMuipData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShareCode   string `protobuf:"bytes,1,opt,name=share_code,json=shareCode,proto3" json:"share_code,omitempty"`
	OwnerUid    uint32 `protobuf:"varint,2,opt,name=owner_uid,json=ownerUid,proto3" json:"owner_uid,omitempty"`
	ModuleId    uint32 `protobuf:"varint,3,opt,name=module_id,json=moduleId,proto3" json:"module_id,omitempty"`
	SceneId     uint32 `protobuf:"varint,4,opt,name=scene_id,json=sceneId,proto3" json:"scene_id,omitempty"`
	BlockId     uint32 `protobuf:"varint,5,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`
	IsAllowCopy bool   `protobuf:"varint,6,opt,name=is_allow_copy,json=isAllowCopy,proto3" json:"is_allow_copy,omitempty"`
	CreateTime  uint32 `protobuf:"varint,7,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
}

func (x *HomeBlueprintBriefMuipData) Reset() {
	*x = HomeBlueprintBriefMuipData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_HomeBlueprintBriefMuipData_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HomeBlueprintBriefMuipData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HomeBlueprintBriefMuipData) ProtoMessage() {}

func (x *HomeBlueprintBriefMuipData) ProtoReflect() protoreflect.Message {
	mi := &file_HomeBlueprintBriefMuipData_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HomeBlueprintBriefMuipData.ProtoReflect.Descriptor instead.
func (*HomeBlueprintBriefMuipData) Descriptor() ([]byte, []int) {
	return file_HomeBlueprintBriefMuipData_proto_rawDescGZIP(), []int{0}
}

func (x *HomeBlueprintBriefMuipData) GetShareCode() string {
	if x != nil {
		return x.ShareCode
	}
	return ""
}

func (x *HomeBlueprintBriefMuipData) GetOwnerUid() uint32 {
	if x != nil {
		return x.OwnerUid
	}
	return 0
}

func (x *HomeBlueprintBriefMuipData) GetModuleId() uint32 {
	if x != nil {
		return x.ModuleId
	}
	return 0
}

func (x *HomeBlueprintBriefMuipData) GetSceneId() uint32 {
	if x != nil {
		return x.SceneId
	}
	return 0
}

func (x *HomeBlueprintBriefMuipData) GetBlockId() uint32 {
	if x != nil {
		return x.BlockId
	}
	return 0
}

func (x *HomeBlueprintBriefMuipData) GetIsAllowCopy() bool {
	if x != nil {
		return x.IsAllowCopy
	}
	return false
}

func (x *HomeBlueprintBriefMuipData) GetCreateTime() uint32 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

var File_HomeBlueprintBriefMuipData_proto protoreflect.FileDescriptor

var file_HomeBlueprintBriefMuipData_proto_rawDesc = []byte{
	0x0a, 0x20, 0x48, 0x6f, 0x6d, 0x65, 0x42, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x42,
	0x72, 0x69, 0x65, 0x66, 0x4d, 0x75, 0x69, 0x70, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf0, 0x01, 0x0a, 0x1a, 0x48, 0x6f,
	0x6d, 0x65, 0x42, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x42, 0x72, 0x69, 0x65, 0x66,
	0x4d, 0x75, 0x69, 0x70, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x77, 0x6e, 0x65, 0x72,
	0x5f, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x55, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49,
	0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x63, 0x65, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x73, 0x63, 0x65, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x61, 0x6c,
	0x6c, 0x6f, 0x77, 0x5f, 0x63, 0x6f, 0x70, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x69, 0x73, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x70, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x0a, 0x5a, 0x08,
	0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_HomeBlueprintBriefMuipData_proto_rawDescOnce sync.Once
	file_HomeBlueprintBriefMuipData_proto_rawDescData = file_HomeBlueprintBriefMuipData_proto_rawDesc
)

func file_HomeBlueprintBriefMuipData_proto_rawDescGZIP() []byte {
	file_HomeBlueprintBriefMuipData_proto_rawDescOnce.Do(func() {
		file_HomeBlueprintBriefMuipData_proto_rawDescData = protoimpl.X.CompressGZIP(file_HomeBlueprintBriefMuipData_proto_rawDescData)
	})
	return file_HomeBlueprintBriefMuipData_proto_rawDescData
}

var file_HomeBlueprintBriefMuipData_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_HomeBlueprintBriefMuipData_proto_goTypes = []interface{}{
	(*HomeBlueprintBriefMuipData)(nil), // 0: proto.HomeBlueprintBriefMuipData
}
var file_HomeBlueprintBriefMuipData_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_HomeBlueprintBriefMuipData_proto_init() }
func file_HomeBlueprintBriefMuipData_proto_init() {
	if File_HomeBlueprintBriefMuipData_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_HomeBlueprintBriefMuipData_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HomeBlueprintBriefMuipData); i {
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
			RawDescriptor: file_HomeBlueprintBriefMuipData_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_HomeBlueprintBriefMuipData_proto_goTypes,
		DependencyIndexes: file_HomeBlueprintBriefMuipData_proto_depIdxs,
		MessageInfos:      file_HomeBlueprintBriefMuipData_proto_msgTypes,
	}.Build()
	File_HomeBlueprintBriefMuipData_proto = out.File
	file_HomeBlueprintBriefMuipData_proto_rawDesc = nil
	file_HomeBlueprintBriefMuipData_proto_goTypes = nil
	file_HomeBlueprintBriefMuipData_proto_depIdxs = nil
}