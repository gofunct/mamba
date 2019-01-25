// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: source/common/source_context.proto

package common

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

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

// `SourceContext` represents information about the source of a
// protobuf element, like the file in which it is defined.
type SourceContext struct {
	// The path-qualified name of the .proto file that contained the associated
	// protobuf element.  For example: `"common/source_context.proto"`.
	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
}

func (m *SourceContext) Reset()         { *m = SourceContext{} }
func (m *SourceContext) String() string { return proto.CompactTextString(m) }
func (*SourceContext) ProtoMessage()    {}
func (*SourceContext) Descriptor() ([]byte, []int) {
	return fileDescriptor_source_context_c9c0fd22410b1d13, []int{0}
}
func (m *SourceContext) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SourceContext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SourceContext.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *SourceContext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SourceContext.Merge(dst, src)
}
func (m *SourceContext) XXX_Size() int {
	return m.Size()
}
func (m *SourceContext) XXX_DiscardUnknown() {
	xxx_messageInfo_SourceContext.DiscardUnknown(m)
}

var xxx_messageInfo_SourceContext proto.InternalMessageInfo

func (m *SourceContext) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func init() {
	proto.RegisterType((*SourceContext)(nil), "common.SourceContext")
}
func (m *SourceContext) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SourceContext) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.FileName) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintSourceContext(dAtA, i, uint64(len(m.FileName)))
		i += copy(dAtA[i:], m.FileName)
	}
	return i, nil
}

func encodeVarintSourceContext(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *SourceContext) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FileName)
	if l > 0 {
		n += 1 + l + sovSourceContext(uint64(l))
	}
	return n
}

func sovSourceContext(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSourceContext(x uint64) (n int) {
	return sovSourceContext(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SourceContext) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSourceContext
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
			return fmt.Errorf("proto: SourceContext: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SourceContext: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FileName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSourceContext
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
				return ErrInvalidLengthSourceContext
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FileName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSourceContext(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSourceContext
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
func skipSourceContext(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSourceContext
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
					return 0, ErrIntOverflowSourceContext
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
					return 0, ErrIntOverflowSourceContext
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
				return 0, ErrInvalidLengthSourceContext
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSourceContext
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
				next, err := skipSourceContext(dAtA[start:])
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
	ErrInvalidLengthSourceContext = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSourceContext   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("source/common/source_context.proto", fileDescriptor_source_context_c9c0fd22410b1d13)
}

var fileDescriptor_source_context_c9c0fd22410b1d13 = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2a, 0xce, 0x2f, 0x2d,
	0x4a, 0x4e, 0xd5, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xd3, 0x87, 0xf0, 0xe2, 0x93, 0xf3, 0xf3,
	0x4a, 0x52, 0x2b, 0x4a, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xd8, 0x20, 0x92, 0x4a, 0x3a,
	0x5c, 0xbc, 0xc1, 0x60, 0x79, 0x67, 0x88, 0xb4, 0x90, 0x34, 0x17, 0x67, 0x5a, 0x66, 0x4e, 0x6a,
	0x7c, 0x5e, 0x62, 0x6e, 0xaa, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x07, 0x48, 0xc0, 0x2f,
	0x31, 0x37, 0xd5, 0x29, 0xf1, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92,
	0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xb8, 0xb8,
	0x92, 0xf3, 0x73, 0xf5, 0x20, 0x66, 0x3a, 0x09, 0xa1, 0x98, 0x18, 0x00, 0xb2, 0x2f, 0x80, 0x71,
	0x11, 0x13, 0xb3, 0x7b, 0x80, 0xd3, 0x2a, 0x26, 0x39, 0xf7, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0x3d,
	0xb0, 0x70, 0x52, 0x69, 0x9a, 0x5e, 0x78, 0x6a, 0x4e, 0x8e, 0x77, 0x5e, 0x7e, 0x79, 0x5e, 0x48,
	0x65, 0x41, 0x6a, 0x71, 0x12, 0x1b, 0xd8, 0x7d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9a,
	0x2e, 0xfd, 0xab, 0xc5, 0x00, 0x00, 0x00,
}
