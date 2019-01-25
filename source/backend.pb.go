// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: source/backend.proto

package source

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import encoding_binary "encoding/binary"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Path Translation specifies how to combine the backend address with the
// request path in order to produce the appropriate forwarding URL for the
// request.
//
// Path Translation is applicable only to HTTP-based backends. Backends which
// do not accept requests over HTTP/HTTPS should leave `path_translation`
// unspecified.
type BackendRule_PathTranslation int32

const (
	BackendRule_PATH_TRANSLATION_UNSPECIFIED BackendRule_PathTranslation = 0
	// Use the backend address as-is, with no modification to the path. If the
	// URL pattern contains variables, the variable names and values will be
	// appended to the query string. If a query string parameter and a URL
	// pattern variable have the same name, this may result in duplicate keys in
	// the query string.
	//
	// # Examples
	//
	// Given the following operation config:
	//
	//     Method path:        /api/company/{cid}/user/{uid}
	//     Backend address:    https://example.cloudfunctions.net/getUser
	//
	// Requests to the following request paths will call the backend at the
	// translated path:
	//
	//     Request path: /api/company/widgetworks/user/johndoe
	//     Translated:   https://example.cloudfunctions.net/getUser?cid=widgetworks&uid=johndoe
	//
	//     Request path: /api/company/widgetworks/user/johndoe?timezone=EST
	//     Translated:   https://example.cloudfunctions.net/getUser?timezone=EST&cid=widgetworks&uid=johndoe
	BackendRule_CONSTANT_ADDRESS BackendRule_PathTranslation = 1
	// The request path will be appended to the backend address.
	//
	// # Examples
	//
	// Given the following operation config:
	//
	//     Method path:        /api/company/{cid}/user/{uid}
	//     Backend address:    https://example.appspot.com
	//
	// Requests to the following request paths will call the backend at the
	// translated path:
	//
	//     Request path: /api/company/widgetworks/user/johndoe
	//     Translated:   https://example.appspot.com/api/company/widgetworks/user/johndoe
	//
	//     Request path: /api/company/widgetworks/user/johndoe?timezone=EST
	//     Translated:   https://example.appspot.com/api/company/widgetworks/user/johndoe?timezone=EST
	BackendRule_APPEND_PATH_TO_ADDRESS BackendRule_PathTranslation = 2
)

var BackendRule_PathTranslation_name = map[int32]string{
	0: "PATH_TRANSLATION_UNSPECIFIED",
	1: "CONSTANT_ADDRESS",
	2: "APPEND_PATH_TO_ADDRESS",
}
var BackendRule_PathTranslation_value = map[string]int32{
	"PATH_TRANSLATION_UNSPECIFIED": 0,
	"CONSTANT_ADDRESS":             1,
	"APPEND_PATH_TO_ADDRESS":       2,
}

func (x BackendRule_PathTranslation) String() string {
	return proto.EnumName(BackendRule_PathTranslation_name, int32(x))
}
func (BackendRule_PathTranslation) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_backend_9ad323b0bb6e6436, []int{1, 0}
}

// `Backend` defines the backend configuration for a service.
type Backend struct {
	// A list of API backend rules that apply to individual API methods.
	//
	// **NOTE:** All service configuration rules follow "last one wins" order.
	Rules []*BackendRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (m *Backend) Reset()         { *m = Backend{} }
func (m *Backend) String() string { return proto.CompactTextString(m) }
func (*Backend) ProtoMessage()    {}
func (*Backend) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_9ad323b0bb6e6436, []int{0}
}
func (m *Backend) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backend) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backend.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Backend) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backend.Merge(dst, src)
}
func (m *Backend) XXX_Size() int {
	return m.Size()
}
func (m *Backend) XXX_DiscardUnknown() {
	xxx_messageInfo_Backend.DiscardUnknown(m)
}

var xxx_messageInfo_Backend proto.InternalMessageInfo

func (m *Backend) GetRules() []*BackendRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

// A backend rule provides configuration for an individual API element.
type BackendRule struct {
	// Selects the methods to which this rule applies.
	//
	// Refer to [selector][google.api.DocumentationRule.selector] for syntax details.
	Selector string `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	// The address of the API backend.
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// The number of seconds to wait for a response from a request.  The default
	// deadline for gRPC is infinite (no deadline) and HTTP requests is 5 seconds.
	Deadline float64 `protobuf:"fixed64,3,opt,name=deadline,proto3" json:"deadline,omitempty"`
	// Minimum deadline in seconds needed for this method. Calls having deadline
	// value lower than this will be rejected.
	MinDeadline float64 `protobuf:"fixed64,4,opt,name=min_deadline,json=minDeadline,proto3" json:"min_deadline,omitempty"`
	// The number of seconds to wait for the completion of a long running
	// operation. The default is no deadline.
	OperationDeadline float64                     `protobuf:"fixed64,5,opt,name=operation_deadline,json=operationDeadline,proto3" json:"operation_deadline,omitempty"`
	PathTranslation   BackendRule_PathTranslation `protobuf:"varint,6,opt,name=path_translation,json=pathTranslation,proto3,enum=source.BackendRule_PathTranslation" json:"path_translation,omitempty"`
	// Authentication settings used by the backend.
	//
	// These are typically used to provide service management functionality to
	// a backend served on a publicly-routable URL. The `authentication`
	// details should match the authentication behavior used by the backend.
	//
	// For example, specifying `jwt_audience` implies that the backend expects
	// authentication via a JWT.
	//
	// Types that are valid to be assigned to Authentication:
	//	*BackendRule_JwtAudience
	Authentication isBackendRule_Authentication `protobuf_oneof:"authentication"`
}

func (m *BackendRule) Reset()         { *m = BackendRule{} }
func (m *BackendRule) String() string { return proto.CompactTextString(m) }
func (*BackendRule) ProtoMessage()    {}
func (*BackendRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_9ad323b0bb6e6436, []int{1}
}
func (m *BackendRule) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BackendRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BackendRule.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *BackendRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackendRule.Merge(dst, src)
}
func (m *BackendRule) XXX_Size() int {
	return m.Size()
}
func (m *BackendRule) XXX_DiscardUnknown() {
	xxx_messageInfo_BackendRule.DiscardUnknown(m)
}

var xxx_messageInfo_BackendRule proto.InternalMessageInfo

type isBackendRule_Authentication interface {
	isBackendRule_Authentication()
	MarshalTo([]byte) (int, error)
	Size() int
}

type BackendRule_JwtAudience struct {
	JwtAudience string `protobuf:"bytes,7,opt,name=jwt_audience,json=jwtAudience,proto3,oneof"`
}

func (*BackendRule_JwtAudience) isBackendRule_Authentication() {}

func (m *BackendRule) GetAuthentication() isBackendRule_Authentication {
	if m != nil {
		return m.Authentication
	}
	return nil
}

func (m *BackendRule) GetSelector() string {
	if m != nil {
		return m.Selector
	}
	return ""
}

func (m *BackendRule) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *BackendRule) GetDeadline() float64 {
	if m != nil {
		return m.Deadline
	}
	return 0
}

func (m *BackendRule) GetMinDeadline() float64 {
	if m != nil {
		return m.MinDeadline
	}
	return 0
}

func (m *BackendRule) GetOperationDeadline() float64 {
	if m != nil {
		return m.OperationDeadline
	}
	return 0
}

func (m *BackendRule) GetPathTranslation() BackendRule_PathTranslation {
	if m != nil {
		return m.PathTranslation
	}
	return BackendRule_PATH_TRANSLATION_UNSPECIFIED
}

func (m *BackendRule) GetJwtAudience() string {
	if x, ok := m.GetAuthentication().(*BackendRule_JwtAudience); ok {
		return x.JwtAudience
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*BackendRule) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _BackendRule_OneofMarshaler, _BackendRule_OneofUnmarshaler, _BackendRule_OneofSizer, []interface{}{
		(*BackendRule_JwtAudience)(nil),
	}
}

func _BackendRule_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*BackendRule)
	// authentication
	switch x := m.Authentication.(type) {
	case *BackendRule_JwtAudience:
		_ = b.EncodeVarint(7<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.JwtAudience)
	case nil:
	default:
		return fmt.Errorf("BackendRule.Authentication has unexpected type %T", x)
	}
	return nil
}

func _BackendRule_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*BackendRule)
	switch tag {
	case 7: // authentication.jwt_audience
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Authentication = &BackendRule_JwtAudience{x}
		return true, err
	default:
		return false, nil
	}
}

func _BackendRule_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*BackendRule)
	// authentication
	switch x := m.Authentication.(type) {
	case *BackendRule_JwtAudience:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.JwtAudience)))
		n += len(x.JwtAudience)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*Backend)(nil), "source.Backend")
	proto.RegisterType((*BackendRule)(nil), "source.BackendRule")
	proto.RegisterEnum("source.BackendRule_PathTranslation", BackendRule_PathTranslation_name, BackendRule_PathTranslation_value)
}
func (m *Backend) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backend) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Rules) > 0 {
		for _, msg := range m.Rules {
			dAtA[i] = 0xa
			i++
			i = encodeVarintBackend(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *BackendRule) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BackendRule) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Selector) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintBackend(dAtA, i, uint64(len(m.Selector)))
		i += copy(dAtA[i:], m.Selector)
	}
	if len(m.Address) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintBackend(dAtA, i, uint64(len(m.Address)))
		i += copy(dAtA[i:], m.Address)
	}
	if m.Deadline != 0 {
		dAtA[i] = 0x19
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Deadline))))
		i += 8
	}
	if m.MinDeadline != 0 {
		dAtA[i] = 0x21
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.MinDeadline))))
		i += 8
	}
	if m.OperationDeadline != 0 {
		dAtA[i] = 0x29
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.OperationDeadline))))
		i += 8
	}
	if m.PathTranslation != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintBackend(dAtA, i, uint64(m.PathTranslation))
	}
	if m.Authentication != nil {
		nn1, err := m.Authentication.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	return i, nil
}

func (m *BackendRule_JwtAudience) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0x3a
	i++
	i = encodeVarintBackend(dAtA, i, uint64(len(m.JwtAudience)))
	i += copy(dAtA[i:], m.JwtAudience)
	return i, nil
}
func encodeVarintBackend(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Backend) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Rules) > 0 {
		for _, e := range m.Rules {
			l = e.Size()
			n += 1 + l + sovBackend(uint64(l))
		}
	}
	return n
}

func (m *BackendRule) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Selector)
	if l > 0 {
		n += 1 + l + sovBackend(uint64(l))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovBackend(uint64(l))
	}
	if m.Deadline != 0 {
		n += 9
	}
	if m.MinDeadline != 0 {
		n += 9
	}
	if m.OperationDeadline != 0 {
		n += 9
	}
	if m.PathTranslation != 0 {
		n += 1 + sovBackend(uint64(m.PathTranslation))
	}
	if m.Authentication != nil {
		n += m.Authentication.Size()
	}
	return n
}

func (m *BackendRule_JwtAudience) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.JwtAudience)
	n += 1 + l + sovBackend(uint64(l))
	return n
}

func sovBackend(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozBackend(x uint64) (n int) {
	return sovBackend(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Backend) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBackend
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Backend: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Backend: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rules", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBackend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBackend
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rules = append(m.Rules, &BackendRule{})
			if err := m.Rules[len(m.Rules)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBackend(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBackend
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BackendRule) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBackend
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BackendRule: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BackendRule: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Selector", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBackend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBackend
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Selector = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBackend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBackend
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deadline", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Deadline = float64(math.Float64frombits(v))
		case 4:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinDeadline", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.MinDeadline = float64(math.Float64frombits(v))
		case 5:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field OperationDeadline", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.OperationDeadline = float64(math.Float64frombits(v))
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PathTranslation", wireType)
			}
			m.PathTranslation = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBackend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PathTranslation |= (BackendRule_PathTranslation(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JwtAudience", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBackend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBackend
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authentication = &BackendRule_JwtAudience{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBackend(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBackend
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBackend(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBackend
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
					return 0, ErrIntOverflowBackend
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBackend
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthBackend
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowBackend
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipBackend(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthBackend = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBackend   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("source/backend.proto", fileDescriptor_backend_9ad323b0bb6e6436) }

var fileDescriptor_backend_9ad323b0bb6e6436 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x41, 0x6e, 0xd3, 0x40,
	0x14, 0x86, 0x3d, 0x49, 0x9b, 0xc0, 0x4b, 0x94, 0x9a, 0xa1, 0x42, 0x56, 0x85, 0x2c, 0x93, 0x6e,
	0xc2, 0x02, 0x23, 0x15, 0xf6, 0x68, 0x5c, 0x1b, 0x6a, 0x09, 0x39, 0xa3, 0xb1, 0x59, 0x5b, 0x53,
	0x7b, 0xd4, 0xb8, 0x38, 0x1e, 0xcb, 0x1e, 0xab, 0xd7, 0xe0, 0x06, 0x48, 0x9c, 0x86, 0x65, 0x97,
	0x2c, 0x51, 0x72, 0x11, 0x14, 0x3b, 0x35, 0x01, 0x75, 0xf9, 0xfe, 0xef, 0x7b, 0x33, 0xfa, 0xa5,
	0x07, 0xa7, 0xb5, 0x6c, 0xaa, 0x44, 0xbc, 0xbd, 0xe6, 0xc9, 0x57, 0x51, 0xa4, 0x76, 0x59, 0x49,
	0x25, 0xf1, 0xa8, 0x4b, 0xe7, 0xef, 0x61, 0xec, 0x74, 0x00, 0xbf, 0x86, 0xe3, 0xaa, 0xc9, 0x45,
	0x6d, 0x20, 0x6b, 0xb8, 0x98, 0x5c, 0x3c, 0xb7, 0x3b, 0xc5, 0xde, 0x73, 0xd6, 0xe4, 0x82, 0x75,
	0xc6, 0xfc, 0xfb, 0x10, 0x26, 0x07, 0x31, 0x3e, 0x83, 0x27, 0xb5, 0xc8, 0x45, 0xa2, 0x64, 0x65,
	0x20, 0x0b, 0x2d, 0x9e, 0xb2, 0x7e, 0xc6, 0x06, 0x8c, 0x79, 0x9a, 0x56, 0xa2, 0xae, 0x8d, 0x41,
	0x8b, 0x1e, 0xc6, 0xdd, 0x56, 0x2a, 0x78, 0x9a, 0x67, 0x85, 0x30, 0x86, 0x16, 0x5a, 0x20, 0xd6,
	0xcf, 0xf8, 0x15, 0x4c, 0xd7, 0x59, 0x11, 0xf7, 0xfc, 0xa8, 0xe5, 0x93, 0x75, 0x56, 0xb8, 0x0f,
	0xca, 0x1b, 0xc0, 0xb2, 0x14, 0x15, 0x57, 0x99, 0x3c, 0x10, 0x8f, 0x5b, 0xf1, 0x59, 0x4f, 0x7a,
	0x3d, 0x00, 0xbd, 0xe4, 0x6a, 0x15, 0xab, 0x8a, 0x17, 0x75, 0xde, 0x32, 0x63, 0x64, 0xa1, 0xc5,
	0xec, 0xe2, 0xfc, 0x91, 0xa6, 0x36, 0xe5, 0x6a, 0x15, 0xfd, 0x55, 0xd9, 0x49, 0xf9, 0x6f, 0x80,
	0xcf, 0x61, 0x7a, 0x7b, 0xa7, 0x62, 0xde, 0xa4, 0x99, 0x28, 0x12, 0x61, 0x8c, 0x77, 0xe5, 0xae,
	0x34, 0x36, 0xb9, 0xbd, 0x53, 0x64, 0x1f, 0xce, 0x05, 0x9c, 0xfc, 0xf7, 0x10, 0xb6, 0xe0, 0x25,
	0x25, 0xd1, 0x55, 0x1c, 0x31, 0x12, 0x84, 0x9f, 0x49, 0xe4, 0x2f, 0x83, 0xf8, 0x4b, 0x10, 0x52,
	0xef, 0xd2, 0xff, 0xe8, 0x7b, 0xae, 0xae, 0xe1, 0x53, 0xd0, 0x2f, 0x97, 0x41, 0x18, 0x91, 0x20,
	0x8a, 0x89, 0xeb, 0x32, 0x2f, 0x0c, 0x75, 0x84, 0xcf, 0xe0, 0x05, 0xa1, 0xd4, 0x0b, 0xdc, 0xb8,
	0x5b, 0x5f, 0xf6, 0x6c, 0xe0, 0xe8, 0x30, 0xe3, 0x8d, 0x5a, 0x89, 0x42, 0x65, 0x49, 0xfb, 0x8b,
	0xf3, 0xe1, 0xe7, 0xc6, 0x44, 0xf7, 0x1b, 0x13, 0xfd, 0xde, 0x98, 0xe8, 0xdb, 0xd6, 0xd4, 0xee,
	0xb7, 0xa6, 0xf6, 0x6b, 0x6b, 0x6a, 0x30, 0x4b, 0xe4, 0xda, 0xbe, 0x91, 0xf2, 0x26, 0x17, 0x36,
	0x2f, 0x33, 0x67, 0xba, 0x6f, 0x4d, 0x77, 0x77, 0x41, 0xd1, 0x8f, 0xc1, 0xd1, 0x27, 0x42, 0xfd,
	0xeb, 0x51, 0x7b, 0x27, 0xef, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0x8b, 0xc3, 0x71, 0xca, 0x3f,
	0x02, 0x00, 0x00,
}