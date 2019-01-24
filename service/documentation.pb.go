// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: service/documentation.proto

package service // import "google.golang.org/genproto/googleapis/api/serviceconfig"

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

// `Documentation` provides the information for describing a service.
//
// Example:
// <pre><code>documentation:
//   summary: >
//     The Google Calendar API gives access
//     to most calendar features.
//   pages:
//   - name: Overview
//     content: &#40;== include google/foo/overview.md ==&#41;
//   - name: Tutorial
//     content: &#40;== include google/foo/tutorial.md ==&#41;
//     subpages;
//     - name: Java
//       content: &#40;== include google/foo/tutorial_java.md ==&#41;
//   rules:
//   - selector: google.calendar.Calendar.Get
//     description: >
//       ...
//   - selector: google.calendar.Calendar.Put
//     description: >
//       ...
// </code></pre>
// Documentation is provided in markdown syntax. In addition to
// standard markdown features, definition lists, tables and fenced
// code blocks are supported. Section headers can be provided and are
// interpreted relative to the section nesting of the context where
// a documentation fragment is embedded.
//
// Documentation from the IDL is merged with documentation defined
// via the config at normalization time, where documentation provided
// by config rules overrides IDL provided.
//
// A number of constructs specific to the API platform are supported
// in documentation text.
//
// In order to reference a proto element, the following
// notation can be used:
// <pre><code>&#91;fully.qualified.proto.name]&#91;]</code></pre>
// To override the display text used for the link, this can be used:
// <pre><code>&#91;display text]&#91;fully.qualified.proto.name]</code></pre>
// Text can be excluded from doc using the following notation:
// <pre><code>&#40;-- internal comment --&#41;</code></pre>
//
// A few directives are available in documentation. Note that
// directives must appear on a single line to be properly
// identified. The `include` directive includes a markdown file from
// an external source:
// <pre><code>&#40;== include path/to/file ==&#41;</code></pre>
// The `resource_for` directive marks a message to be the resource of
// a collection in REST view. If it is not specified, tools attempt
// to infer the resource from the operations in a collection:
// <pre><code>&#40;== resource_for v1.shelves.books ==&#41;</code></pre>
// The directive `suppress_warning` does not directly affect documentation
// and is documented together with service config validation.
type Documentation struct {
	// A short summary of what the service does. Can only be provided by
	// plain text.
	Summary string `protobuf:"bytes,1,opt,name=summary,proto3" json:"summary,omitempty"`
	// The top level pages for the documentation set.
	Pages []*Page `protobuf:"bytes,5,rep,name=pages,proto3" json:"pages,omitempty"`
	// A list of documentation rules that apply to individual API elements.
	//
	// **NOTE:** All service configuration rules follow "last one wins" order.
	Rules []*DocumentationRule `protobuf:"bytes,3,rep,name=rules,proto3" json:"rules,omitempty"`
	// The URL to the root of documentation.
	DocumentationRootUrl string `protobuf:"bytes,4,opt,name=documentation_root_url,json=documentationRootUrl,proto3" json:"documentation_root_url,omitempty"`
	// Declares a single overview page. For example:
	// <pre><code>documentation:
	//   summary: ...
	//   overview: &#40;== include overview.md ==&#41;
	// </code></pre>
	// This is a shortcut for the following declaration (using pages style):
	// <pre><code>documentation:
	//   summary: ...
	//   pages:
	//   - name: Overview
	//     content: &#40;== include overview.md ==&#41;
	// </code></pre>
	// Note: you cannot specify both `overview` field and `pages` field.
	Overview string `protobuf:"bytes,2,opt,name=overview,proto3" json:"overview,omitempty"`
}

func (m *Documentation) Reset()         { *m = Documentation{} }
func (m *Documentation) String() string { return proto.CompactTextString(m) }
func (*Documentation) ProtoMessage()    {}
func (*Documentation) Descriptor() ([]byte, []int) {
	return fileDescriptor_documentation_e61e0c1dbd00e8d3, []int{0}
}
func (m *Documentation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Documentation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Documentation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Documentation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Documentation.Merge(dst, src)
}
func (m *Documentation) XXX_Size() int {
	return m.Size()
}
func (m *Documentation) XXX_DiscardUnknown() {
	xxx_messageInfo_Documentation.DiscardUnknown(m)
}

var xxx_messageInfo_Documentation proto.InternalMessageInfo

func (m *Documentation) GetSummary() string {
	if m != nil {
		return m.Summary
	}
	return ""
}

func (m *Documentation) GetPages() []*Page {
	if m != nil {
		return m.Pages
	}
	return nil
}

func (m *Documentation) GetRules() []*DocumentationRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

func (m *Documentation) GetDocumentationRootUrl() string {
	if m != nil {
		return m.DocumentationRootUrl
	}
	return ""
}

func (m *Documentation) GetOverview() string {
	if m != nil {
		return m.Overview
	}
	return ""
}

// A documentation rule provides information about individual API elements.
type DocumentationRule struct {
	// The selector is a comma-separated list of patterns. Each pattern is a
	// qualified name of the element which may end in "*", indicating a wildcard.
	// Wildcards are only allowed at the end and for a whole component of the
	// qualified name, i.e. "foo.*" is ok, but not "foo.b*" or "foo.*.bar". To
	// specify a default for all applicable elements, the whole pattern "*"
	// is used.
	Selector string `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	// Description of the selected API(s).
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Deprecation description of the selected element(s). It can be provided if an
	// element is marked as `deprecated`.
	DeprecationDescription string `protobuf:"bytes,3,opt,name=deprecation_description,json=deprecationDescription,proto3" json:"deprecation_description,omitempty"`
}

func (m *DocumentationRule) Reset()         { *m = DocumentationRule{} }
func (m *DocumentationRule) String() string { return proto.CompactTextString(m) }
func (*DocumentationRule) ProtoMessage()    {}
func (*DocumentationRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_documentation_e61e0c1dbd00e8d3, []int{1}
}
func (m *DocumentationRule) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DocumentationRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DocumentationRule.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *DocumentationRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DocumentationRule.Merge(dst, src)
}
func (m *DocumentationRule) XXX_Size() int {
	return m.Size()
}
func (m *DocumentationRule) XXX_DiscardUnknown() {
	xxx_messageInfo_DocumentationRule.DiscardUnknown(m)
}

var xxx_messageInfo_DocumentationRule proto.InternalMessageInfo

func (m *DocumentationRule) GetSelector() string {
	if m != nil {
		return m.Selector
	}
	return ""
}

func (m *DocumentationRule) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *DocumentationRule) GetDeprecationDescription() string {
	if m != nil {
		return m.DeprecationDescription
	}
	return ""
}

// Represents a documentation page. A page can contain subpages to represent
// nested documentation set structure.
type Page struct {
	// The name of the page. It will be used as an identity of the page to
	// generate URI of the page, text of the link to this page in navigation,
	// etc. The full page name (start from the root page name to this page
	// concatenated with `.`) can be used as reference to the page in your
	// documentation. For example:
	// <pre><code>pages:
	// - name: Tutorial
	//   content: &#40;== include tutorial.md ==&#41;
	//   subpages:
	//   - name: Java
	//     content: &#40;== include tutorial_java.md ==&#41;
	// </code></pre>
	// You can reference `Java` page using Markdown reference link syntax:
	// `[Java][Tutorial.Java]`.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The Markdown content of the page. You can use <code>&#40;== include {path} ==&#41;</code>
	// to include content from a Markdown file.
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	// Subpages of this page. The order of subpages specified here will be
	// honored in the generated docset.
	Subpages []*Page `protobuf:"bytes,3,rep,name=subpages,proto3" json:"subpages,omitempty"`
}

func (m *Page) Reset()         { *m = Page{} }
func (m *Page) String() string { return proto.CompactTextString(m) }
func (*Page) ProtoMessage()    {}
func (*Page) Descriptor() ([]byte, []int) {
	return fileDescriptor_documentation_e61e0c1dbd00e8d3, []int{2}
}
func (m *Page) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Page) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Page.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Page) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Page.Merge(dst, src)
}
func (m *Page) XXX_Size() int {
	return m.Size()
}
func (m *Page) XXX_DiscardUnknown() {
	xxx_messageInfo_Page.DiscardUnknown(m)
}

var xxx_messageInfo_Page proto.InternalMessageInfo

func (m *Page) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Page) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Page) GetSubpages() []*Page {
	if m != nil {
		return m.Subpages
	}
	return nil
}

func init() {
	proto.RegisterType((*Documentation)(nil), "service.Documentation")
	proto.RegisterType((*DocumentationRule)(nil), "service.DocumentationRule")
	proto.RegisterType((*Page)(nil), "service.Page")
}
func (m *Documentation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Documentation) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Summary) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.Summary)))
		i += copy(dAtA[i:], m.Summary)
	}
	if len(m.Overview) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.Overview)))
		i += copy(dAtA[i:], m.Overview)
	}
	if len(m.Rules) > 0 {
		for _, msg := range m.Rules {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintDocumentation(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.DocumentationRootUrl) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.DocumentationRootUrl)))
		i += copy(dAtA[i:], m.DocumentationRootUrl)
	}
	if len(m.Pages) > 0 {
		for _, msg := range m.Pages {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintDocumentation(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *DocumentationRule) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DocumentationRule) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Selector) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.Selector)))
		i += copy(dAtA[i:], m.Selector)
	}
	if len(m.Description) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.Description)))
		i += copy(dAtA[i:], m.Description)
	}
	if len(m.DeprecationDescription) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.DeprecationDescription)))
		i += copy(dAtA[i:], m.DeprecationDescription)
	}
	return i, nil
}

func (m *Page) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Page) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Content) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDocumentation(dAtA, i, uint64(len(m.Content)))
		i += copy(dAtA[i:], m.Content)
	}
	if len(m.Subpages) > 0 {
		for _, msg := range m.Subpages {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintDocumentation(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintDocumentation(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Documentation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Summary)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	l = len(m.Overview)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	if len(m.Rules) > 0 {
		for _, e := range m.Rules {
			l = e.Size()
			n += 1 + l + sovDocumentation(uint64(l))
		}
	}
	l = len(m.DocumentationRootUrl)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	if len(m.Pages) > 0 {
		for _, e := range m.Pages {
			l = e.Size()
			n += 1 + l + sovDocumentation(uint64(l))
		}
	}
	return n
}

func (m *DocumentationRule) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Selector)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	l = len(m.DeprecationDescription)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	return n
}

func (m *Page) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	l = len(m.Content)
	if l > 0 {
		n += 1 + l + sovDocumentation(uint64(l))
	}
	if len(m.Subpages) > 0 {
		for _, e := range m.Subpages {
			l = e.Size()
			n += 1 + l + sovDocumentation(uint64(l))
		}
	}
	return n
}

func sovDocumentation(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozDocumentation(x uint64) (n int) {
	return sovDocumentation(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Documentation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDocumentation
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
			return fmt.Errorf("proto: Documentation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Documentation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Summary", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Summary = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Overview", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Overview = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rules", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rules = append(m.Rules, &DocumentationRule{})
			if err := m.Rules[len(m.Rules)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DocumentationRootUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DocumentationRootUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pages = append(m.Pages, &Page{})
			if err := m.Pages[len(m.Pages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDocumentation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDocumentation
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
func (m *DocumentationRule) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDocumentation
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
			return fmt.Errorf("proto: DocumentationRule: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DocumentationRule: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Selector", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Selector = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeprecationDescription", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeprecationDescription = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDocumentation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDocumentation
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
func (m *Page) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDocumentation
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
			return fmt.Errorf("proto: Page: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Page: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Content = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subpages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDocumentation
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
				return ErrInvalidLengthDocumentation
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subpages = append(m.Subpages, &Page{})
			if err := m.Subpages[len(m.Subpages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDocumentation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDocumentation
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
func skipDocumentation(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDocumentation
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
					return 0, ErrIntOverflowDocumentation
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
					return 0, ErrIntOverflowDocumentation
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
				return 0, ErrInvalidLengthDocumentation
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDocumentation
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
				next, err := skipDocumentation(dAtA[start:])
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
	ErrInvalidLengthDocumentation = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDocumentation   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("service/documentation.proto", fileDescriptor_documentation_e61e0c1dbd00e8d3)
}

var fileDescriptor_documentation_e61e0c1dbd00e8d3 = []byte{
	// 384 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xc1, 0xca, 0xd3, 0x40,
	0x14, 0x85, 0x3b, 0x7f, 0x52, 0xff, 0x3a, 0xa5, 0x82, 0x83, 0xd4, 0x50, 0x21, 0x94, 0xba, 0xa9,
	0x9b, 0x44, 0x54, 0x70, 0xe1, 0xca, 0x52, 0x11, 0x77, 0x21, 0xe0, 0xc6, 0x4d, 0x99, 0x4e, 0xaf,
	0x43, 0x20, 0x99, 0x1b, 0x66, 0x26, 0x15, 0x5f, 0x41, 0x5c, 0xf8, 0x0c, 0x3e, 0x8d, 0xcb, 0xe2,
	0xca, 0xa5, 0xb4, 0x2f, 0x22, 0x49, 0xa6, 0x31, 0x41, 0x77, 0x39, 0xf9, 0xce, 0xcc, 0x39, 0x77,
	0xb8, 0xf4, 0x91, 0x01, 0x7d, 0xcc, 0x04, 0xc4, 0x07, 0x14, 0x55, 0x01, 0xca, 0x72, 0x9b, 0xa1,
	0x8a, 0x4a, 0x8d, 0x16, 0xd9, 0xad, 0x83, 0xab, 0x9f, 0x84, 0xce, 0xb6, 0x7d, 0x03, 0x0b, 0xe8,
	0xad, 0xa9, 0x8a, 0x82, 0xeb, 0xcf, 0x01, 0x59, 0x92, 0xf5, 0xdd, 0xf4, 0x2a, 0xd9, 0x82, 0x4e,
	0xf0, 0x58, 0x9f, 0x83, 0x4f, 0xc1, 0x4d, 0x83, 0x3a, 0xcd, 0x9e, 0xd2, 0xb1, 0xae, 0x72, 0x30,
	0x81, 0xb7, 0xf4, 0xd6, 0xd3, 0x67, 0x8b, 0xc8, 0x05, 0x44, 0x83, 0xcb, 0xd3, 0x2a, 0x87, 0xb4,
	0x35, 0xb2, 0x17, 0x74, 0x3e, 0x68, 0xb6, 0xd3, 0x88, 0x76, 0x57, 0xe9, 0x3c, 0xf0, 0x9b, 0xbb,
	0x1f, 0x0c, 0x68, 0x8a, 0x68, 0xdf, 0xeb, 0x9c, 0x3d, 0xa6, 0xe3, 0x92, 0x4b, 0x30, 0xc1, 0xb8,
	0xc9, 0x99, 0x75, 0x39, 0x09, 0x97, 0x90, 0xb6, 0x6c, 0xf5, 0x85, 0xd0, 0xfb, 0xff, 0xe4, 0xd6,
	0xf5, 0x0d, 0xe4, 0x20, 0x2c, 0x6a, 0x37, 0x59, 0xa7, 0xd9, 0x92, 0x4e, 0x0f, 0x60, 0x84, 0xce,
	0xca, 0xda, 0xee, 0xa6, 0xeb, 0xff, 0x62, 0x2f, 0xe9, 0xc3, 0x03, 0x94, 0x1a, 0x44, 0x5b, 0xb6,
	0xef, 0xf6, 0x1a, 0xf7, 0xbc, 0x87, 0xb7, 0x7f, 0xe9, 0x6a, 0x47, 0xfd, 0xba, 0x1b, 0x63, 0xd4,
	0x57, 0xbc, 0x00, 0x17, 0xdd, 0x7c, 0xd7, 0x6f, 0x2d, 0x50, 0x59, 0x50, 0xd6, 0x45, 0x5e, 0x25,
	0x7b, 0x42, 0x27, 0xa6, 0xda, 0xb7, 0xa3, 0x7a, 0xff, 0x1b, 0xb5, 0xc3, 0x9b, 0xaf, 0xe4, 0xc7,
	0x39, 0x24, 0xa7, 0x73, 0x48, 0x7e, 0x9f, 0x43, 0xf2, 0xed, 0x12, 0x8e, 0x4e, 0x97, 0x70, 0xf4,
	0xeb, 0x12, 0x8e, 0xe8, 0x3d, 0x81, 0x45, 0x24, 0x11, 0x65, 0x0e, 0x11, 0x2f, 0xb3, 0x0d, 0x1b,
	0xbc, 0x4a, 0x52, 0xaf, 0x42, 0x42, 0x3e, 0xbc, 0x71, 0x0e, 0x89, 0x39, 0x57, 0x32, 0x42, 0x2d,
	0x63, 0x09, 0xaa, 0x59, 0x94, 0xb8, 0x45, 0xbc, 0xcc, 0x4c, 0xcc, 0xcb, 0x2c, 0x76, 0x15, 0x04,
	0xaa, 0x8f, 0x99, 0x7c, 0x35, 0x50, 0xdf, 0x6f, 0xfc, 0xb7, 0xaf, 0x93, 0x77, 0xfb, 0x3b, 0xcd,
	0xc1, 0xe7, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x61, 0x80, 0xf7, 0x18, 0x80, 0x02, 0x00, 0x00,
}