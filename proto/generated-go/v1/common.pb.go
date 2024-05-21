// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: v1/common.proto

package v1

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

type State int32

const (
	State_STATE_UNSPECIFIED State = 0
	State_ACTIVE            State = 1
	State_DELETED           State = 2
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "STATE_UNSPECIFIED",
		1: "ACTIVE",
		2: "DELETED",
	}
	State_value = map[string]int32{
		"STATE_UNSPECIFIED": 0,
		"ACTIVE":            1,
		"DELETED":           2,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_common_proto_enumTypes[0].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_v1_common_proto_enumTypes[0]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_v1_common_proto_rawDescGZIP(), []int{0}
}

type Engine int32

const (
	Engine_ENGINE_UNSPECIFIED Engine = 0
	Engine_CLICKHOUSE         Engine = 1
	Engine_MYSQL              Engine = 2
	Engine_POSTGRES           Engine = 3
	Engine_SNOWFLAKE          Engine = 4
	Engine_SQLITE             Engine = 5
	Engine_TIDB               Engine = 6
	Engine_MONGODB            Engine = 7
	Engine_REDIS              Engine = 8
	Engine_ORACLE             Engine = 9
	Engine_SPANNER            Engine = 10
	Engine_MSSQL              Engine = 11
	Engine_REDSHIFT           Engine = 12
	Engine_MARIADB            Engine = 13
	Engine_OCEANBASE          Engine = 14
	Engine_DM                 Engine = 15
	Engine_RISINGWAVE         Engine = 16
	Engine_OCEANBASE_ORACLE   Engine = 17
	Engine_STARROCKS          Engine = 18
	Engine_DORIS              Engine = 19
	Engine_HIVE               Engine = 20
	Engine_ELASTICSEARCH      Engine = 21
)

// Enum value maps for Engine.
var (
	Engine_name = map[int32]string{
		0:  "ENGINE_UNSPECIFIED",
		1:  "CLICKHOUSE",
		2:  "MYSQL",
		3:  "POSTGRES",
		4:  "SNOWFLAKE",
		5:  "SQLITE",
		6:  "TIDB",
		7:  "MONGODB",
		8:  "REDIS",
		9:  "ORACLE",
		10: "SPANNER",
		11: "MSSQL",
		12: "REDSHIFT",
		13: "MARIADB",
		14: "OCEANBASE",
		15: "DM",
		16: "RISINGWAVE",
		17: "OCEANBASE_ORACLE",
		18: "STARROCKS",
		19: "DORIS",
		20: "HIVE",
		21: "ELASTICSEARCH",
	}
	Engine_value = map[string]int32{
		"ENGINE_UNSPECIFIED": 0,
		"CLICKHOUSE":         1,
		"MYSQL":              2,
		"POSTGRES":           3,
		"SNOWFLAKE":          4,
		"SQLITE":             5,
		"TIDB":               6,
		"MONGODB":            7,
		"REDIS":              8,
		"ORACLE":             9,
		"SPANNER":            10,
		"MSSQL":              11,
		"REDSHIFT":           12,
		"MARIADB":            13,
		"OCEANBASE":          14,
		"DM":                 15,
		"RISINGWAVE":         16,
		"OCEANBASE_ORACLE":   17,
		"STARROCKS":          18,
		"DORIS":              19,
		"HIVE":               20,
		"ELASTICSEARCH":      21,
	}
)

func (x Engine) Enum() *Engine {
	p := new(Engine)
	*p = x
	return p
}

func (x Engine) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Engine) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_common_proto_enumTypes[1].Descriptor()
}

func (Engine) Type() protoreflect.EnumType {
	return &file_v1_common_proto_enumTypes[1]
}

func (x Engine) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Engine.Descriptor instead.
func (Engine) EnumDescriptor() ([]byte, []int) {
	return file_v1_common_proto_rawDescGZIP(), []int{1}
}

type VCSType int32

const (
	VCSType_VCS_TYPE_UNSPECIFIED VCSType = 0
	// GitHub type. Using for GitHub community edition(ce).
	VCSType_GITHUB VCSType = 1
	// GitLab type. Using for GitLab community edition(ce) and enterprise edition(ee).
	VCSType_GITLAB VCSType = 2
	// BitBucket type. Using for BitBucket cloud or BitBucket server.
	VCSType_BITBUCKET VCSType = 3
	// Azure DevOps. Using for Azure DevOps GitOps workflow.
	VCSType_AZURE_DEVOPS VCSType = 4
)

// Enum value maps for VCSType.
var (
	VCSType_name = map[int32]string{
		0: "VCS_TYPE_UNSPECIFIED",
		1: "GITHUB",
		2: "GITLAB",
		3: "BITBUCKET",
		4: "AZURE_DEVOPS",
	}
	VCSType_value = map[string]int32{
		"VCS_TYPE_UNSPECIFIED": 0,
		"GITHUB":               1,
		"GITLAB":               2,
		"BITBUCKET":            3,
		"AZURE_DEVOPS":         4,
	}
)

func (x VCSType) Enum() *VCSType {
	p := new(VCSType)
	*p = x
	return p
}

func (x VCSType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VCSType) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_common_proto_enumTypes[2].Descriptor()
}

func (VCSType) Type() protoreflect.EnumType {
	return &file_v1_common_proto_enumTypes[2]
}

func (x VCSType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VCSType.Descriptor instead.
func (VCSType) EnumDescriptor() ([]byte, []int) {
	return file_v1_common_proto_rawDescGZIP(), []int{2}
}

type MaskingLevel int32

const (
	MaskingLevel_MASKING_LEVEL_UNSPECIFIED MaskingLevel = 0
	MaskingLevel_NONE                      MaskingLevel = 1
	MaskingLevel_PARTIAL                   MaskingLevel = 2
	MaskingLevel_FULL                      MaskingLevel = 3
)

// Enum value maps for MaskingLevel.
var (
	MaskingLevel_name = map[int32]string{
		0: "MASKING_LEVEL_UNSPECIFIED",
		1: "NONE",
		2: "PARTIAL",
		3: "FULL",
	}
	MaskingLevel_value = map[string]int32{
		"MASKING_LEVEL_UNSPECIFIED": 0,
		"NONE":                      1,
		"PARTIAL":                   2,
		"FULL":                      3,
	}
)

func (x MaskingLevel) Enum() *MaskingLevel {
	p := new(MaskingLevel)
	*p = x
	return p
}

func (x MaskingLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MaskingLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_common_proto_enumTypes[3].Descriptor()
}

func (MaskingLevel) Type() protoreflect.EnumType {
	return &file_v1_common_proto_enumTypes[3]
}

func (x MaskingLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MaskingLevel.Descriptor instead.
func (MaskingLevel) EnumDescriptor() ([]byte, []int) {
	return file_v1_common_proto_rawDescGZIP(), []int{3}
}

type ExportFormat int32

const (
	ExportFormat_FORMAT_UNSPECIFIED ExportFormat = 0
	ExportFormat_CSV                ExportFormat = 1
	ExportFormat_JSON               ExportFormat = 2
	ExportFormat_SQL                ExportFormat = 3
	ExportFormat_XLSX               ExportFormat = 4
)

// Enum value maps for ExportFormat.
var (
	ExportFormat_name = map[int32]string{
		0: "FORMAT_UNSPECIFIED",
		1: "CSV",
		2: "JSON",
		3: "SQL",
		4: "XLSX",
	}
	ExportFormat_value = map[string]int32{
		"FORMAT_UNSPECIFIED": 0,
		"CSV":                1,
		"JSON":               2,
		"SQL":                3,
		"XLSX":               4,
	}
)

func (x ExportFormat) Enum() *ExportFormat {
	p := new(ExportFormat)
	*p = x
	return p
}

func (x ExportFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ExportFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_common_proto_enumTypes[4].Descriptor()
}

func (ExportFormat) Type() protoreflect.EnumType {
	return &file_v1_common_proto_enumTypes[4]
}

func (x ExportFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ExportFormat.Descriptor instead.
func (ExportFormat) EnumDescriptor() ([]byte, []int) {
	return file_v1_common_proto_rawDescGZIP(), []int{4}
}

var File_v1_common_proto protoreflect.FileDescriptor

var file_v1_common_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0b, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2a, 0x37,
	0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x53, 0x54, 0x41, 0x54, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45,
	0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x02, 0x2a, 0xb9, 0x02, 0x0a, 0x06, 0x45, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x4e, 0x47, 0x49, 0x4e, 0x45, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4c,
	0x49, 0x43, 0x4b, 0x48, 0x4f, 0x55, 0x53, 0x45, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x59,
	0x53, 0x51, 0x4c, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x4f, 0x53, 0x54, 0x47, 0x52, 0x45,
	0x53, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x4e, 0x4f, 0x57, 0x46, 0x4c, 0x41, 0x4b, 0x45,
	0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x51, 0x4c, 0x49, 0x54, 0x45, 0x10, 0x05, 0x12, 0x08,
	0x0a, 0x04, 0x54, 0x49, 0x44, 0x42, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x4f, 0x4e, 0x47,
	0x4f, 0x44, 0x42, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x45, 0x44, 0x49, 0x53, 0x10, 0x08,
	0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x52, 0x41, 0x43, 0x4c, 0x45, 0x10, 0x09, 0x12, 0x0b, 0x0a, 0x07,
	0x53, 0x50, 0x41, 0x4e, 0x4e, 0x45, 0x52, 0x10, 0x0a, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x53, 0x53,
	0x51, 0x4c, 0x10, 0x0b, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x44, 0x53, 0x48, 0x49, 0x46, 0x54,
	0x10, 0x0c, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x41, 0x52, 0x49, 0x41, 0x44, 0x42, 0x10, 0x0d, 0x12,
	0x0d, 0x0a, 0x09, 0x4f, 0x43, 0x45, 0x41, 0x4e, 0x42, 0x41, 0x53, 0x45, 0x10, 0x0e, 0x12, 0x06,
	0x0a, 0x02, 0x44, 0x4d, 0x10, 0x0f, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x49, 0x53, 0x49, 0x4e, 0x47,
	0x57, 0x41, 0x56, 0x45, 0x10, 0x10, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x43, 0x45, 0x41, 0x4e, 0x42,
	0x41, 0x53, 0x45, 0x5f, 0x4f, 0x52, 0x41, 0x43, 0x4c, 0x45, 0x10, 0x11, 0x12, 0x0d, 0x0a, 0x09,
	0x53, 0x54, 0x41, 0x52, 0x52, 0x4f, 0x43, 0x4b, 0x53, 0x10, 0x12, 0x12, 0x09, 0x0a, 0x05, 0x44,
	0x4f, 0x52, 0x49, 0x53, 0x10, 0x13, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x49, 0x56, 0x45, 0x10, 0x14,
	0x12, 0x11, 0x0a, 0x0d, 0x45, 0x4c, 0x41, 0x53, 0x54, 0x49, 0x43, 0x53, 0x45, 0x41, 0x52, 0x43,
	0x48, 0x10, 0x15, 0x2a, 0x5c, 0x0a, 0x07, 0x56, 0x43, 0x53, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18,
	0x0a, 0x14, 0x56, 0x43, 0x53, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x49, 0x54, 0x48,
	0x55, 0x42, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x49, 0x54, 0x4c, 0x41, 0x42, 0x10, 0x02,
	0x12, 0x0d, 0x0a, 0x09, 0x42, 0x49, 0x54, 0x42, 0x55, 0x43, 0x4b, 0x45, 0x54, 0x10, 0x03, 0x12,
	0x10, 0x0a, 0x0c, 0x41, 0x5a, 0x55, 0x52, 0x45, 0x5f, 0x44, 0x45, 0x56, 0x4f, 0x50, 0x53, 0x10,
	0x04, 0x2a, 0x4e, 0x0a, 0x0c, 0x4d, 0x61, 0x73, 0x6b, 0x69, 0x6e, 0x67, 0x4c, 0x65, 0x76, 0x65,
	0x6c, 0x12, 0x1d, 0x0a, 0x19, 0x4d, 0x41, 0x53, 0x4b, 0x49, 0x4e, 0x47, 0x5f, 0x4c, 0x45, 0x56,
	0x45, 0x4c, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x41,
	0x52, 0x54, 0x49, 0x41, 0x4c, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x55, 0x4c, 0x4c, 0x10,
	0x03, 0x2a, 0x4c, 0x0a, 0x0c, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x12, 0x16, 0x0a, 0x12, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x43, 0x53, 0x56,
	0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4a, 0x53, 0x4f, 0x4e, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03,
	0x53, 0x51, 0x4c, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x58, 0x4c, 0x53, 0x58, 0x10, 0x04, 0x42,
	0x11, 0x5a, 0x0f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2d, 0x67, 0x6f, 0x2f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1_common_proto_rawDescOnce sync.Once
	file_v1_common_proto_rawDescData = file_v1_common_proto_rawDesc
)

func file_v1_common_proto_rawDescGZIP() []byte {
	file_v1_common_proto_rawDescOnce.Do(func() {
		file_v1_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_common_proto_rawDescData)
	})
	return file_v1_common_proto_rawDescData
}

var file_v1_common_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_v1_common_proto_goTypes = []interface{}{
	(State)(0),        // 0: bytebase.v1.State
	(Engine)(0),       // 1: bytebase.v1.Engine
	(VCSType)(0),      // 2: bytebase.v1.VCSType
	(MaskingLevel)(0), // 3: bytebase.v1.MaskingLevel
	(ExportFormat)(0), // 4: bytebase.v1.ExportFormat
}
var file_v1_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_common_proto_init() }
func file_v1_common_proto_init() {
	if File_v1_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1_common_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_common_proto_goTypes,
		DependencyIndexes: file_v1_common_proto_depIdxs,
		EnumInfos:         file_v1_common_proto_enumTypes,
	}.Build()
	File_v1_common_proto = out.File
	file_v1_common_proto_rawDesc = nil
	file_v1_common_proto_goTypes = nil
	file_v1_common_proto_depIdxs = nil
}
