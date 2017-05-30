// Code generated by protoc-gen-go.
// source: proto/warpdrive.proto
// DO NOT EDIT!

/*
Package warpdrive is a generated protocol buffer package.

It is generated from these files:
	proto/warpdrive.proto

It has these top-level messages:
	Release
	Chunck
*/
package warpdrive

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Platform int32

const (
	Platform_UNKNOWN Platform = 0
	Platform_IOS     Platform = 1
	Platform_ANDROID Platform = 2
)

var Platform_name = map[int32]string{
	0: "UNKNOWN",
	1: "IOS",
	2: "ANDROID",
}
var Platform_value = map[string]int32{
	"UNKNOWN": 0,
	"IOS":     1,
	"ANDROID": 2,
}

func (x Platform) String() string {
	return proto.EnumName(Platform_name, int32(x))
}
func (Platform) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Release can be duplicate for different rollout.
// for example, bundle can be released under beta for App1 and then
// it needs to be pushed again to be rollout as production. behind the scene both records are
// using the same bundle binary but they are tragetign two sets of people.
type Release struct {
	// @inject_tag: storm:"id,increment"
	Id uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty" storm:"id,increment"`
	// @inject_tag: storm:"index"
	App string `protobuf:"bytes,2,opt,name=app" json:"app,omitempty" storm:"index"`
	// this is just for label. it's not unique
	// becuase you might want to rollback
	// @inject_tag: storm:"index"
	Version  string   `protobuf:"bytes,3,opt,name=version" json:"version,omitempty" storm:"index"`
	Notes    string   `protobuf:"bytes,4,opt,name=notes" json:"notes,omitempty"`
	Platform Platform `protobuf:"varint,5,opt,name=platform,enum=warpdrive.Platform" json:"platform,omitempty"`
	// this is list of releases that can safely upgrade to this
	// version.
	NextReleaseId uint64 `protobuf:"varint,6,opt,name=nextReleaseId" json:"nextReleaseId,omitempty" `
	// this is used as what kind of release is. As an example `beta`
	RolloutAt string `protobuf:"bytes,7,opt,name=rolloutAt" json:"rolloutAt,omitempty" `
	// this is the hash value of bundle package
	// @inject_tag: storm:"index"
	Bundle string `protobuf:"bytes,8,opt,name=bundle" json:"bundle,omitempty" storm:"index"`
	// if the lock value is true, it means that this release can not be ultered or modified.
	// this is used to make sure the production doesn't download the unlock one.
	Lock bool `protobuf:"varint,9,opt,name=lock" json:"lock,omitempty" `
	// @inject_tag: storm:"index"
	CreatedAt string `protobuf:"bytes,10,opt,name=createdAt" json:"createdAt,omitempty" storm:"index"`
	// @inject_tag: storm:"index"
	UpdatedAt string `protobuf:"bytes,11,opt,name=updatedAt" json:"updatedAt,omitempty" storm:"index"`
}

func (m *Release) Reset()                    { *m = Release{} }
func (m *Release) String() string            { return proto.CompactTextString(m) }
func (*Release) ProtoMessage()               {}
func (*Release) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Release) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Release) GetApp() string {
	if m != nil {
		return m.App
	}
	return ""
}

func (m *Release) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Release) GetNotes() string {
	if m != nil {
		return m.Notes
	}
	return ""
}

func (m *Release) GetPlatform() Platform {
	if m != nil {
		return m.Platform
	}
	return Platform_UNKNOWN
}

func (m *Release) GetNextReleaseId() uint64 {
	if m != nil {
		return m.NextReleaseId
	}
	return 0
}

func (m *Release) GetRolloutAt() string {
	if m != nil {
		return m.RolloutAt
	}
	return ""
}

func (m *Release) GetBundle() string {
	if m != nil {
		return m.Bundle
	}
	return ""
}

func (m *Release) GetLock() bool {
	if m != nil {
		return m.Lock
	}
	return false
}

func (m *Release) GetCreatedAt() string {
	if m != nil {
		return m.CreatedAt
	}
	return ""
}

func (m *Release) GetUpdatedAt() string {
	if m != nil {
		return m.UpdatedAt
	}
	return ""
}

type Chunck struct {
	// Types that are valid to be assigned to Value:
	//	*Chunck_Header_
	//	*Chunck_Body_
	Value isChunck_Value `protobuf_oneof:"value" `
}

func (m *Chunck) Reset()                    { *m = Chunck{} }
func (m *Chunck) String() string            { return proto.CompactTextString(m) }
func (*Chunck) ProtoMessage()               {}
func (*Chunck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isChunck_Value interface {
	isChunck_Value()
}

type Chunck_Header_ struct {
	Header *Chunck_Header `protobuf:"bytes,1,opt,name=header,oneof"`
}
type Chunck_Body_ struct {
	Body *Chunck_Body `protobuf:"bytes,2,opt,name=body,oneof"`
}

func (*Chunck_Header_) isChunck_Value() {}
func (*Chunck_Body_) isChunck_Value()   {}

func (m *Chunck) GetValue() isChunck_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Chunck) GetHeader() *Chunck_Header {
	if x, ok := m.GetValue().(*Chunck_Header_); ok {
		return x.Header
	}
	return nil
}

func (m *Chunck) GetBody() *Chunck_Body {
	if x, ok := m.GetValue().(*Chunck_Body_); ok {
		return x.Body
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Chunck) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Chunck_OneofMarshaler, _Chunck_OneofUnmarshaler, _Chunck_OneofSizer, []interface{}{
		(*Chunck_Header_)(nil),
		(*Chunck_Body_)(nil),
	}
}

func _Chunck_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Chunck)
	// value
	switch x := m.Value.(type) {
	case *Chunck_Header_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Header); err != nil {
			return err
		}
	case *Chunck_Body_:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Body); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Chunck.Value has unexpected type %T", x)
	}
	return nil
}

func _Chunck_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Chunck)
	switch tag {
	case 1: // value.header
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Chunck_Header)
		err := b.DecodeMessage(msg)
		m.Value = &Chunck_Header_{msg}
		return true, err
	case 2: // value.body
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Chunck_Body)
		err := b.DecodeMessage(msg)
		m.Value = &Chunck_Body_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Chunck_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Chunck)
	// value
	switch x := m.Value.(type) {
	case *Chunck_Header_:
		s := proto.Size(x.Header)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Chunck_Body_:
		s := proto.Size(x.Body)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Chunck_Header struct {
	ReleaseId uint64 `protobuf:"varint,1,opt,name=releaseId" json:"releaseId,omitempty"`
}

func (m *Chunck_Header) Reset()                    { *m = Chunck_Header{} }
func (m *Chunck_Header) String() string            { return proto.CompactTextString(m) }
func (*Chunck_Header) ProtoMessage()               {}
func (*Chunck_Header) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *Chunck_Header) GetReleaseId() uint64 {
	if m != nil {
		return m.ReleaseId
	}
	return 0
}

type Chunck_Body struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Chunck_Body) Reset()                    { *m = Chunck_Body{} }
func (m *Chunck_Body) String() string            { return proto.CompactTextString(m) }
func (*Chunck_Body) ProtoMessage()               {}
func (*Chunck_Body) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 1} }

func (m *Chunck_Body) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Release)(nil), "warpdrive.Release")
	proto.RegisterType((*Chunck)(nil), "warpdrive.Chunck")
	proto.RegisterType((*Chunck_Header)(nil), "warpdrive.Chunck.Header")
	proto.RegisterType((*Chunck_Body)(nil), "warpdrive.Chunck.Body")
	proto.RegisterEnum("warpdrive.Platform", Platform_name, Platform_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Command service

type CommandClient interface {
	CreateRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Release, error)
	// if Release.Id is provided, then only the matched one returns.
	// if Release.App is provided, then it returns all the releases for that app
	GetRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (Command_GetReleaseClient, error)
	// once the release.lock set to true, Release can not be updated anymore,
	// only `nextReleaseId` can be changed under the following condition:
	// nextReleaseId must not set or `lock` has to be false
	UpdateRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Release, error)
	// UplaodRelease won't work unless ReleaseId exists
	UploadRelease(ctx context.Context, opts ...grpc.CallOption) (Command_UploadReleaseClient, error)
}

type commandClient struct {
	cc *grpc.ClientConn
}

func NewCommandClient(cc *grpc.ClientConn) CommandClient {
	return &commandClient{cc}
}

func (c *commandClient) CreateRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Release, error) {
	out := new(Release)
	err := grpc.Invoke(ctx, "/warpdrive.Command/CreateRelease", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandClient) GetRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (Command_GetReleaseClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Command_serviceDesc.Streams[0], c.cc, "/warpdrive.Command/GetRelease", opts...)
	if err != nil {
		return nil, err
	}
	x := &commandGetReleaseClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Command_GetReleaseClient interface {
	Recv() (*Release, error)
	grpc.ClientStream
}

type commandGetReleaseClient struct {
	grpc.ClientStream
}

func (x *commandGetReleaseClient) Recv() (*Release, error) {
	m := new(Release)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *commandClient) UpdateRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Release, error) {
	out := new(Release)
	err := grpc.Invoke(ctx, "/warpdrive.Command/UpdateRelease", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandClient) UploadRelease(ctx context.Context, opts ...grpc.CallOption) (Command_UploadReleaseClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Command_serviceDesc.Streams[1], c.cc, "/warpdrive.Command/UploadRelease", opts...)
	if err != nil {
		return nil, err
	}
	x := &commandUploadReleaseClient{stream}
	return x, nil
}

type Command_UploadReleaseClient interface {
	Send(*Chunck) error
	CloseAndRecv() (*Release, error)
	grpc.ClientStream
}

type commandUploadReleaseClient struct {
	grpc.ClientStream
}

func (x *commandUploadReleaseClient) Send(m *Chunck) error {
	return x.ClientStream.SendMsg(m)
}

func (x *commandUploadReleaseClient) CloseAndRecv() (*Release, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Release)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Command service

type CommandServer interface {
	CreateRelease(context.Context, *Release) (*Release, error)
	// if Release.Id is provided, then only the matched one returns.
	// if Release.App is provided, then it returns all the releases for that app
	GetRelease(*Release, Command_GetReleaseServer) error
	// once the release.lock set to true, Release can not be updated anymore,
	// only `nextReleaseId` can be changed under the following condition:
	// nextReleaseId must not set or `lock` has to be false
	UpdateRelease(context.Context, *Release) (*Release, error)
	// UplaodRelease won't work unless ReleaseId exists
	UploadRelease(Command_UploadReleaseServer) error
}

func RegisterCommandServer(s *grpc.Server, srv CommandServer) {
	s.RegisterService(&_Command_serviceDesc, srv)
}

func _Command_CreateRelease_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServer).CreateRelease(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/warpdrive.Command/CreateRelease",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServer).CreateRelease(ctx, req.(*Release))
	}
	return interceptor(ctx, in, info, handler)
}

func _Command_GetRelease_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Release)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CommandServer).GetRelease(m, &commandGetReleaseServer{stream})
}

type Command_GetReleaseServer interface {
	Send(*Release) error
	grpc.ServerStream
}

type commandGetReleaseServer struct {
	grpc.ServerStream
}

func (x *commandGetReleaseServer) Send(m *Release) error {
	return x.ServerStream.SendMsg(m)
}

func _Command_UpdateRelease_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServer).UpdateRelease(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/warpdrive.Command/UpdateRelease",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServer).UpdateRelease(ctx, req.(*Release))
	}
	return interceptor(ctx, in, info, handler)
}

func _Command_UploadRelease_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CommandServer).UploadRelease(&commandUploadReleaseServer{stream})
}

type Command_UploadReleaseServer interface {
	SendAndClose(*Release) error
	Recv() (*Chunck, error)
	grpc.ServerStream
}

type commandUploadReleaseServer struct {
	grpc.ServerStream
}

func (x *commandUploadReleaseServer) SendAndClose(m *Release) error {
	return x.ServerStream.SendMsg(m)
}

func (x *commandUploadReleaseServer) Recv() (*Chunck, error) {
	m := new(Chunck)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Command_serviceDesc = grpc.ServiceDesc{
	ServiceName: "warpdrive.Command",
	HandlerType: (*CommandServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRelease",
			Handler:    _Command_CreateRelease_Handler,
		},
		{
			MethodName: "UpdateRelease",
			Handler:    _Command_UpdateRelease_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetRelease",
			Handler:       _Command_GetRelease_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadRelease",
			Handler:       _Command_UploadRelease_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/warpdrive.proto",
}

// Client API for Query service

type QueryClient interface {
	// the folowing four fields must be presented in Release object
	// `Release.id`, `Release.app`, `Release.platform`, `Release.rolloutAt`
	// when client need to know the next Release
	GetUpgrade(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Release, error)
	DownloadRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (Query_DownloadReleaseClient, error)
}

type queryClient struct {
	cc *grpc.ClientConn
}

func NewQueryClient(cc *grpc.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GetUpgrade(ctx context.Context, in *Release, opts ...grpc.CallOption) (*Release, error) {
	out := new(Release)
	err := grpc.Invoke(ctx, "/warpdrive.Query/GetUpgrade", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DownloadRelease(ctx context.Context, in *Release, opts ...grpc.CallOption) (Query_DownloadReleaseClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Query_serviceDesc.Streams[0], c.cc, "/warpdrive.Query/DownloadRelease", opts...)
	if err != nil {
		return nil, err
	}
	x := &queryDownloadReleaseClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Query_DownloadReleaseClient interface {
	Recv() (*Chunck, error)
	grpc.ClientStream
}

type queryDownloadReleaseClient struct {
	grpc.ClientStream
}

func (x *queryDownloadReleaseClient) Recv() (*Chunck, error) {
	m := new(Chunck)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Query service

type QueryServer interface {
	// the folowing four fields must be presented in Release object
	// `Release.id`, `Release.app`, `Release.platform`, `Release.rolloutAt`
	// when client need to know the next Release
	GetUpgrade(context.Context, *Release) (*Release, error)
	DownloadRelease(*Release, Query_DownloadReleaseServer) error
}

func RegisterQueryServer(s *grpc.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GetUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/warpdrive.Query/GetUpgrade",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetUpgrade(ctx, req.(*Release))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DownloadRelease_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Release)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(QueryServer).DownloadRelease(m, &queryDownloadReleaseServer{stream})
}

type Query_DownloadReleaseServer interface {
	Send(*Chunck) error
	grpc.ServerStream
}

type queryDownloadReleaseServer struct {
	grpc.ServerStream
}

func (x *queryDownloadReleaseServer) Send(m *Chunck) error {
	return x.ServerStream.SendMsg(m)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "warpdrive.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUpgrade",
			Handler:    _Query_GetUpgrade_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadRelease",
			Handler:       _Query_DownloadRelease_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/warpdrive.proto",
}

func init() { proto.RegisterFile("proto/warpdrive.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 471 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x41, 0x8f, 0xd2, 0x40,
	0x14, 0x66, 0xba, 0xa5, 0x85, 0x87, 0xbb, 0xe2, 0x53, 0x37, 0x13, 0xe2, 0x81, 0x10, 0x63, 0x1a,
	0xa3, 0xac, 0xa9, 0x89, 0x51, 0xe3, 0x85, 0x85, 0x44, 0x88, 0x09, 0x68, 0x0d, 0xf1, 0x3c, 0x30,
	0xa3, 0x4b, 0xb6, 0x74, 0x9a, 0x61, 0xca, 0xca, 0xc1, 0xdf, 0xa5, 0x7f, 0xcc, 0xbb, 0x99, 0x69,
	0x29, 0x18, 0xb8, 0x70, 0x9b, 0xf7, 0x7d, 0xef, 0x7b, 0xef, 0x7b, 0x1f, 0x14, 0x1e, 0xa7, 0x4a,
	0x6a, 0x79, 0x75, 0xc7, 0x54, 0xca, 0xd5, 0x62, 0x2d, 0xba, 0xb6, 0xc6, 0x7a, 0x09, 0x74, 0xfe,
	0x38, 0xe0, 0x47, 0x22, 0x16, 0x6c, 0x25, 0xf0, 0x02, 0x9c, 0x05, 0xa7, 0xa4, 0x4d, 0x02, 0x37,
	0x72, 0x16, 0x1c, 0x9b, 0x70, 0xc6, 0xd2, 0x94, 0x3a, 0x6d, 0x12, 0xd4, 0x23, 0xf3, 0x44, 0x0a,
	0xfe, 0x5a, 0xa8, 0xd5, 0x42, 0x26, 0xf4, 0xcc, 0xa2, 0xdb, 0x12, 0x1f, 0x41, 0x35, 0x91, 0x5a,
	0xac, 0xa8, 0x6b, 0xf1, 0xbc, 0xc0, 0x2b, 0xa8, 0xa5, 0x31, 0xd3, 0xdf, 0xa5, 0x5a, 0xd2, 0x6a,
	0x9b, 0x04, 0x17, 0xe1, 0xc3, 0xee, 0xce, 0xcc, 0xe7, 0x82, 0x8a, 0xca, 0x26, 0x7c, 0x0a, 0xe7,
	0x89, 0xf8, 0xa9, 0x0b, 0x47, 0x23, 0x4e, 0x3d, 0xeb, 0xe6, 0x7f, 0x10, 0x9f, 0x40, 0x5d, 0xc9,
	0x38, 0x96, 0x99, 0xee, 0x69, 0xea, 0xdb, 0x85, 0x3b, 0x00, 0x2f, 0xc1, 0x9b, 0x65, 0x09, 0x8f,
	0x05, 0xad, 0x59, 0xaa, 0xa8, 0x10, 0xc1, 0x8d, 0xe5, 0xfc, 0x96, 0xd6, 0xdb, 0x24, 0xa8, 0x45,
	0xf6, 0x6d, 0x26, 0xcd, 0x95, 0x60, 0x5a, 0xf0, 0x9e, 0xa6, 0x90, 0x4f, 0x2a, 0x01, 0xc3, 0x66,
	0x29, 0x2f, 0xd8, 0x46, 0xce, 0x96, 0x40, 0xe7, 0x37, 0x01, 0xaf, 0x7f, 0x93, 0x25, 0xf3, 0x5b,
	0x0c, 0xc1, 0xbb, 0x11, 0x8c, 0x0b, 0x65, 0xd3, 0x6b, 0x84, 0x74, 0xef, 0xca, 0xbc, 0xa5, 0x3b,
	0xb4, 0xfc, 0xb0, 0x12, 0x15, 0x9d, 0xf8, 0x02, 0xdc, 0x99, 0xe4, 0x1b, 0x1b, 0x6f, 0x23, 0xbc,
	0x3c, 0x54, 0x5c, 0x4b, 0xbe, 0x19, 0x56, 0x22, 0xdb, 0xd5, 0x7a, 0x06, 0x5e, 0x3e, 0xc1, 0x1e,
	0x5f, 0xc6, 0x93, 0xff, 0x58, 0x3b, 0xa0, 0xd5, 0x02, 0xd7, 0xe8, 0xcc, 0xb1, 0x9c, 0x69, 0x66,
	0x1b, 0xee, 0x45, 0xf6, 0x7d, 0xed, 0x43, 0x75, 0xcd, 0xe2, 0x4c, 0x3c, 0x7f, 0x09, 0xb5, 0x6d,
	0xf6, 0xd8, 0x00, 0x7f, 0x3a, 0xfe, 0x34, 0x9e, 0x7c, 0x1b, 0x37, 0x2b, 0xe8, 0xc3, 0xd9, 0x68,
	0xf2, 0xb5, 0x49, 0x0c, 0xda, 0x1b, 0x0f, 0xa2, 0xc9, 0x68, 0xd0, 0x74, 0xc2, 0xbf, 0x04, 0xfc,
	0xbe, 0x5c, 0x2e, 0x59, 0xc2, 0xf1, 0x1d, 0x9c, 0xf7, 0x6d, 0x3e, 0xdb, 0x3f, 0x0d, 0xee, 0x19,
	0x2f, 0xb0, 0xd6, 0x11, 0xac, 0x53, 0xc1, 0xb7, 0x00, 0x1f, 0x85, 0x3e, 0x59, 0xf7, 0x8a, 0x98,
	0xa5, 0x53, 0x1b, 0xfb, 0xe9, 0x4b, 0xdf, 0x1b, 0x69, 0x2c, 0x19, 0xdf, 0x4a, 0x1f, 0x1c, 0x04,
	0x7d, 0x5c, 0x19, 0x90, 0xf0, 0x17, 0x54, 0xbf, 0x64, 0x42, 0x6d, 0xf0, 0x8d, 0x75, 0x3e, 0x4d,
	0x7f, 0x28, 0xc6, 0x4f, 0x59, 0xfe, 0x01, 0xee, 0x0f, 0xe4, 0x5d, 0xb2, 0xbf, 0xfe, 0x98, 0xf8,
	0xd0, 0x92, 0xb9, 0x7a, 0xe6, 0xd9, 0x8f, 0xf5, 0xf5, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x14,
	0x58, 0x11, 0xe6, 0xc5, 0x03, 0x00, 0x00,
}
