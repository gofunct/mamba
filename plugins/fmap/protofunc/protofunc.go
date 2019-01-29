package protofunc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	ggdescriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"strings"
)

func StringMethodOptionsExtension(fieldID int32, f *descriptor.MethodDescriptorProto) string {
	if f == nil {
		return ""
	}
	if f.Options == nil {
		return ""
	}
	var extendedType *descriptor.MethodOptions
	var extensionType *string

	eds := proto.RegisteredExtensions(f.Options)
	if eds[fieldID] == nil {
		ed := &proto.ExtensionDesc{
			ExtendedType:  extendedType,
			ExtensionType: extensionType,
			Field:         fieldID,
			Tag:           fmt.Sprintf("bytes,%d", fieldID),
		}
		proto.RegisterExtension(ed)
		eds = proto.RegisteredExtensions(f.Options)
	}

	ext, err := proto.GetExtension(f.Options, eds[fieldID])
	if err != nil {
		return ""
	}

	str, ok := ext.(*string)
	if !ok {
		return ""
	}

	return *str
}

func StringFieldExtension(fieldID int32, f *descriptor.FieldDescriptorProto) string {
	if f == nil {
		return ""
	}
	if f.Options == nil {
		return ""
	}
	var extendedType *descriptor.FieldOptions
	var extensionType *string

	eds := proto.RegisteredExtensions(f.Options)
	if eds[fieldID] == nil {
		ed := &proto.ExtensionDesc{
			ExtendedType:  extendedType,
			ExtensionType: extensionType,
			Field:         fieldID,
			Tag:           fmt.Sprintf("bytes,%d", fieldID),
		}
		proto.RegisterExtension(ed)
		eds = proto.RegisteredExtensions(f.Options)
	}

	ext, err := proto.GetExtension(f.Options, eds[fieldID])
	if err != nil {
		return ""
	}

	str, ok := ext.(*string)
	if !ok {
		return ""
	}

	return *str
}

func Int64FieldExtension(fieldID int32, f *descriptor.FieldDescriptorProto) int64 {
	if f == nil {
		return 0
	}
	if f.Options == nil {
		return 0
	}
	var extendedType *descriptor.FieldOptions
	var extensionType *string

	eds := proto.RegisteredExtensions(f.Options)
	if eds[fieldID] == nil {
		ed := &proto.ExtensionDesc{
			ExtendedType:  extendedType,
			ExtensionType: extensionType,
			Field:         fieldID,
			Tag:           fmt.Sprintf("bytes,%d", fieldID),
		}
		proto.RegisterExtension(ed)
		eds = proto.RegisteredExtensions(f.Options)
	}

	ext, err := proto.GetExtension(f.Options, eds[fieldID])
	if err != nil {
		return 0
	}

	i, ok := ext.(*int64)
	if !ok {
		return 0
	}

	return *i
}

func BoolMethodOptionsExtension(fieldID int32, f *descriptor.MethodDescriptorProto) bool {
	if f == nil {
		return false
	}
	if f.Options == nil {
		return false
	}
	var extendedType *descriptor.MethodOptions
	var extensionType *bool

	eds := proto.RegisteredExtensions(f.Options)
	if eds[fieldID] == nil {
		ed := &proto.ExtensionDesc{
			ExtendedType:  extendedType,
			ExtensionType: extensionType,
			Field:         fieldID,
			Tag:           fmt.Sprintf("varint,%d", fieldID),
		}
		proto.RegisterExtension(ed)
		eds = proto.RegisteredExtensions(f.Options)
	}

	ext, err := proto.GetExtension(f.Options, eds[fieldID])
	if err != nil {
		return false
	}

	b, ok := ext.(*bool)
	if !ok {
		return false
	}

	return *b
}

func BoolFieldExtension(fieldID int32, f *descriptor.FieldDescriptorProto) bool {
	if f == nil {
		return false
	}
	if f.Options == nil {
		return false
	}
	var extendedType *descriptor.FieldOptions
	var extensionType *bool

	eds := proto.RegisteredExtensions(f.Options)
	if eds[fieldID] == nil {
		ed := &proto.ExtensionDesc{
			ExtendedType:  extendedType,
			ExtensionType: extensionType,
			Field:         fieldID,
			Tag:           fmt.Sprintf("varint,%d", fieldID),
		}
		proto.RegisterExtension(ed)
		eds = proto.RegisteredExtensions(f.Options)
	}

	ext, err := proto.GetExtension(f.Options, eds[fieldID])
	if err != nil {
		return false
	}

	b, ok := ext.(*bool)
	if !ok {
		return false
	}

	return *b
}

func GetProtoFile(name string) *ggdescriptor.File {
	if registry == nil {
		return nil
	}
	file, err := registry.LookupFile(name)
	if err != nil {
		panic(err)
	}
	return file
}

func GetMessageType(f *descriptor.FileDescriptorProto, name string) *ggdescriptor.Message {
	if registry != nil {
		msg, err := registry.LookupMsg(".", name)
		if err != nil {
			panic(err)
		}
		return msg
	}

	// name is in the form .packageName.MessageTypeName.InnerMessageTypeName...
	// e.g. .article.ProductTag
	splits := strings.Split(name, ".")
	target := splits[len(splits)-1]
	for _, m := range f.MessageType {
		if target == *m.Name {
			return &ggdescriptor.Message{
				DescriptorProto: m,
			}
		}
	}
	return nil
}

func GetEnumValue(f []*descriptor.EnumDescriptorProto, name string) []*descriptor.EnumValueDescriptorProto {
	for _, item := range f {
		if strings.EqualFold(*item.Name, name) {
			return item.GetValue()
		}
	}

	return nil
}

func IsFieldMessageTimeStamp(f *descriptor.FieldDescriptorProto) bool {
	if f.Type != nil && *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		if strings.Compare(*f.TypeName, ".google.protobuf.Timestamp") == 0 {
			return true
		}
	}
	return false
}

func IsFieldMessage(f *descriptor.FieldDescriptorProto) bool {
	if f.Type != nil && *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return true
	}

	return false
}

func IsFieldRepeated(f *descriptor.FieldDescriptorProto) bool {
	if f == nil {
		return false
	}
	if f.Type != nil && f.Label != nil && *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		return true
	}

	return false
}

func IsFieldMap(f *descriptor.FieldDescriptorProto, m *descriptor.DescriptorProto) bool {
	if f.TypeName == nil {
		return false
	}

	shortName := ShortType(*f.TypeName)
	var nt *descriptor.DescriptorProto
	for _, t := range m.NestedType {
		if *t.Name == shortName {
			nt = t
			break
		}
	}

	if nt == nil {
		return false
	}

	for _, f := range nt.Field {
		switch *f.Name {
		case "key":
			if *f.Number != 1 {
				return false
			}
		case "value":
			if *f.Number != 2 {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func FieldMapKeyType(f *descriptor.FieldDescriptorProto, m *descriptor.DescriptorProto) *descriptor.FieldDescriptorProto {
	if f.TypeName == nil {
		return nil
	}

	shortName := ShortType(*f.TypeName)
	var nt *descriptor.DescriptorProto
	for _, t := range m.NestedType {
		if *t.Name == shortName {
			nt = t
			break
		}
	}

	if nt == nil {
		return nil
	}

	for _, f := range nt.Field {
		if *f.Name == "key" {
			return f
		}
	}

	return nil

}

func FieldMapValType(f *descriptor.FieldDescriptorProto, m *descriptor.DescriptorProto) *descriptor.FieldDescriptorProto {
	if f.TypeName == nil {
		return nil
	}

	shortName := ShortType(*f.TypeName)
	var nt *descriptor.DescriptorProto
	for _, t := range m.NestedType {
		if *t.Name == shortName {
			nt = t
			break
		}
	}

	if nt == nil {
		return nil
	}

	for _, f := range nt.Field {
		if *f.Name == "value" {
			return f
		}
	}

	return nil

}

func GoTypeWithGoPackage(p *descriptor.FileDescriptorProto, f *descriptor.FieldDescriptorProto) string {
	pkg := ""
	if *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE || *f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM {
		if IsTimestampPackage(*f.TypeName) {
			pkg = "timestamp"
		} else {
			pkg = *p.GetOptions().GoPackage
			if strings.Contains(*p.GetOptions().GoPackage, ";") {
				pkg = strings.Split(*p.GetOptions().GoPackage, ";")[1]
			}
		}
	}
	return GoType(pkg, f, p)
}

func HaskellType(pkg string, f *descriptor.FieldDescriptorProto) string {
	switch *f.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Float]"
		}
		return "Float"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Float]"
		}
		return "Float"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Int64]"
		}
		return "Int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Word]"
		}
		return "Word"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Int]"
		}
		return "Int"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Word]"
		}
		return "Word"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Bool]"
		}
		return "Bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Text]"
		}
		return "Text"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if pkg != "" {
			pkg = pkg + "."
		}
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return fmt.Sprintf("[%s%s]", pkg, ShortType(*f.TypeName))
		}
		return fmt.Sprintf("%s%s", pkg, ShortType(*f.TypeName))
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Word8]"
		}
		return "Word8"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return fmt.Sprintf("%s%s", pkg, ShortType(*f.TypeName))
	default:
		return "Generic"
	}
}

func GoType(pkg string, f *descriptor.FieldDescriptorProto, p *descriptor.FileDescriptorProto) string {
	if pkg != "" {
		pkg = pkg + "."
	}
	switch *f.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]float64"
		}
		return "float64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]float32"
		}
		return "float32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int64"
		}
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint64"
		}
		return "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int32"
		}
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint32"
		}
		return "uint32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]bool"
		}
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]string"
		}
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		name := *f.TypeName
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			fieldPackage := strings.Split(*f.TypeName, ".")
			filePackage := strings.Split(*p.Package, ".")
			// check if we are working with a message embedded.
			if len(fieldPackage) > 1 && len(fieldPackage)+1 > len(filePackage)+1 {
				name = strings.Join(fieldPackage[len(filePackage)+1:], "_")
			}

			return fmt.Sprintf("[]*%s%s", pkg, ShortType(name))
		}
		return fmt.Sprintf("*%s%s", pkg, ShortType(name))
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]byte"
		}
		return "byte"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		name := *f.TypeName
		fieldPackage := strings.Split(*f.TypeName, ".")
		filePackage := strings.Split(*p.Package, ".")
		// check if we are working with a message embedded.
		if len(fieldPackage) > 1 && len(fieldPackage)+1 > len(filePackage)+1 {
			name = strings.Join(fieldPackage[len(filePackage)+1:], "_")
		}
		return fmt.Sprintf("*%s%s", pkg, ShortType(name))
	default:
		return "interface{}"
	}
}

func GoZeroValue(f *descriptor.FieldDescriptorProto) string {
	const nilString = "nil"
	if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		return nilString
	}
	switch *f.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "0.0"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "0.0"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return "0"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "0"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return "0"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return "0"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "false"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "\"\""
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return nilString
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "0"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return nilString
	default:
		return nilString
	}
}

func JsType(f *descriptor.FieldDescriptorProto) string {
	template := "%s"
	if IsFieldRepeated(f) {
		template = "Array<%s>"
	}

	switch *f.Type {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_ENUM:
		return fmt.Sprintf(template, NamespacedFlowType(*f.TypeName))
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return fmt.Sprintf(template, "number")
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return fmt.Sprintf(template, "boolean")
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return fmt.Sprintf(template, "Uint8Array")
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return fmt.Sprintf(template, "string")
	default:
		return fmt.Sprintf(template, "any")
	}
}

func JsSuffixReservedKeyword(s string) string {
	return jsReservedRe.ReplaceAllString(s, "${1}${2}_${3}")
}

func IsTimestampPackage(s string) bool {
	var isTimestampPackage bool
	if strings.Compare(s, ".google.protobuf.Timestamp") == 0 {
		isTimestampPackage = true
	}
	return isTimestampPackage
}

func ShortType(s string) string {
	t := strings.Split(s, ".")
	return t[len(t)-1]
}

func LeadingComment(i interface{}) string {
	loc := PathMap[i]
	return loc.GetLeadingComments()
}
func TrailingComment(i interface{}) string {
	loc := PathMap[i]
	return loc.GetTrailingComments()
}
func LeadingDetachedComments(i interface{}) []string {
	loc := PathMap[i]
	return loc.GetLeadingDetachedComments()
}

func NamespacedFlowType(s string) string {
	trimmed := strings.TrimLeft(s, ".")
	splitted := strings.Split(trimmed, ".")
	return strings.Join(splitted, "$")
}

func SetRegistry(reg *ggdescriptor.Registry) {
	registry = reg
}
