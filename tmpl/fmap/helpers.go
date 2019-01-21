package fmap

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	ggdescriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/huandu/xstrings"
	"regexp"
	"strings"
)

var jsReservedRe = regexp.MustCompile(`(^|[^A-Za-z])(do|if|in|for|let|new|try|var|case|else|enum|eval|false|null|this|true|void|with|break|catch|class|const|super|throw|while|yield|delete|export|import|public|return|static|switch|typeof|default|extends|finally|package|private|continue|debugger|function|arguments|interface|protected|implements|instanceof)($|[^A-Za-z])`)

var (
	registry *ggdescriptor.Registry // some helpers need access to registry
)

var pathMap map[interface{}]*descriptor.SourceCodeInfo_Location

func SetRegistry(reg *ggdescriptor.Registry) {
	registry = reg
}

func InitPathMap(file *descriptor.FileDescriptorProto) {
	pathMap = make(map[interface{}]*descriptor.SourceCodeInfo_Location)
	addToPathMap(file.GetSourceCodeInfo(), file, []int32{})
}

func InitPathMaps(files []*descriptor.FileDescriptorProto) {
	pathMap = make(map[interface{}]*descriptor.SourceCodeInfo_Location)
	for _, file := range files {
		addToPathMap(file.GetSourceCodeInfo(), file, []int32{})
	}
}

func addToPathMap(info *descriptor.SourceCodeInfo, i interface{}, path []int32) {
	loc := findLoc(info, path)
	if loc != nil {
		pathMap[i] = loc
	}
	switch d := i.(type) {
	case *descriptor.FileDescriptorProto:
		for index, descriptor := range d.MessageType {
			addToPathMap(info, descriptor, newPath(path, 4, index))
		}
		for index, descriptor := range d.EnumType {
			addToPathMap(info, descriptor, newPath(path, 5, index))
		}
		for index, descriptor := range d.Service {
			addToPathMap(info, descriptor, newPath(path, 6, index))
		}
	case *descriptor.DescriptorProto:
		for index, descriptor := range d.Field {
			addToPathMap(info, descriptor, newPath(path, 2, index))
		}
		for index, descriptor := range d.NestedType {
			addToPathMap(info, descriptor, newPath(path, 3, index))
		}
		for index, descriptor := range d.EnumType {
			addToPathMap(info, descriptor, newPath(path, 4, index))
		}
	case *descriptor.EnumDescriptorProto:
		for index, descriptor := range d.Value {
			addToPathMap(info, descriptor, newPath(path, 2, index))
		}
	case *descriptor.ServiceDescriptorProto:
		for index, descriptor := range d.Method {
			addToPathMap(info, descriptor, newPath(path, 2, index))
		}
	}
}

func newPath(base []int32, field int32, index int) []int32 {
	p := append([]int32{}, base...)
	p = append(p, field, int32(index))
	return p
}

func findLoc(info *descriptor.SourceCodeInfo, path []int32) *descriptor.SourceCodeInfo_Location {
	for _, loc := range info.GetLocation() {
		if samePath(loc.Path, path) {
			return loc
		}
	}
	return nil
}

func samePath(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, p := range a {
		if p != b[i] {
			return false
		}
	}
	return true
}

func findSourceInfoLocation(i interface{}) *descriptor.SourceCodeInfo_Location {
	if pathMap == nil {
		return nil
	}
	return pathMap[i]
}

func leadingComment(i interface{}) string {
	loc := pathMap[i]
	return loc.GetLeadingComments()
}
func trailingComment(i interface{}) string {
	loc := pathMap[i]
	return loc.GetTrailingComments()
}
func leadingDetachedComments(i interface{}) []string {
	loc := pathMap[i]
	return loc.GetLeadingDetachedComments()
}

func stringMethodOptionsExtension(fieldID int32, f *descriptor.MethodDescriptorProto) string {
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

func stringFieldExtension(fieldID int32, f *descriptor.FieldDescriptorProto) string {
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

func boolMethodOptionsExtension(fieldID int32, f *descriptor.MethodDescriptorProto) bool {
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
			Tag:           fmt.Sprintf("bytes,%d", fieldID),
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

func boolFieldExtension(fieldID int32, f *descriptor.FieldDescriptorProto) bool {
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

func init() {
	for k, v := range sprig.TxtFuncMap() {
		ProtoHelpersFuncMap[k] = v
	}
}

func getProtoFile(name string) *ggdescriptor.File {
	if registry == nil {
		return nil
	}
	file, err := registry.LookupFile(name)
	if err != nil {
		panic(err)
	}
	return file
}

func getMessageType(f *descriptor.FileDescriptorProto, name string) *ggdescriptor.Message {
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

func getEnumValue(f []*descriptor.EnumDescriptorProto, name string) []*descriptor.EnumValueDescriptorProto {
	for _, item := range f {
		if strings.EqualFold(*item.Name, name) {
			return item.GetValue()
		}
	}

	return nil
}

func isFieldMessageTimeStamp(f *descriptor.FieldDescriptorProto) bool {
	if f.Type != nil && *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		if strings.Compare(*f.TypeName, ".google.protobuf.Timestamp") == 0 {
			return true
		}
	}
	return false
}

func isFieldMessage(f *descriptor.FieldDescriptorProto) bool {
	if f.Type != nil && *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return true
	}

	return false
}

func isFieldRepeated(f *descriptor.FieldDescriptorProto) bool {
	if f == nil {
		return false
	}
	if f.Type != nil && f.Label != nil && *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		return true
	}

	return false
}

func isFieldMap(f *descriptor.FieldDescriptorProto, m *descriptor.DescriptorProto) bool {
	if f.TypeName == nil {
		return false
	}

	shortName := shortType(*f.TypeName)
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

func fieldMapKeyType(f *descriptor.FieldDescriptorProto, m *descriptor.DescriptorProto) *descriptor.FieldDescriptorProto {
	if f.TypeName == nil {
		return nil
	}

	shortName := shortType(*f.TypeName)
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

func fieldMapValueType(f *descriptor.FieldDescriptorProto, m *descriptor.DescriptorProto) *descriptor.FieldDescriptorProto {
	if f.TypeName == nil {
		return nil
	}

	shortName := shortType(*f.TypeName)
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

func isTimestampPackage(s string) bool {
	var isTimestampPackage bool
	if strings.Compare(s, ".google.protobuf.Timestamp") == 0 {
		isTimestampPackage = true
	}
	return isTimestampPackage
}

func getPackageTypeName(s string) string {
	if strings.Contains(s, ".") {
		return strings.Split(s, ".")[1]
	}
	return ""
}

func shortType(s string) string {
	t := strings.Split(s, ".")
	return t[len(t)-1]
}

func namespacedFlowType(s string) string {
	trimmed := strings.TrimLeft(s, ".")
	splitted := strings.Split(trimmed, ".")
	return strings.Join(splitted, "$")
}

func lowerGoNormalize(s string) string {
	fmtd := xstrings.ToCamelCase(s)
	fmtd = xstrings.FirstRuneToLower(fmtd)
	return formatID(s, fmtd)
}

func goNormalize(s string) string {
	fmtd := xstrings.ToCamelCase(s)
	return formatID(s, fmtd)
}

func formatID(base string, formatted string) string {
	if formatted == "" {
		return formatted
	}
	switch {
	case base == "id":
		// id -> ID
		return "ID"
	case strings.HasPrefix(base, "id_"):
		// id_some -> IDSome
		return "ID" + formatted[2:]
	case strings.HasSuffix(base, "_id"):
		// some_id -> SomeID
		return formatted[:len(formatted)-2] + "ID"
	case strings.HasSuffix(base, "_ids"):
		// some_ids -> SomeIDs
		return formatted[:len(formatted)-3] + "IDs"
	}
	return formatted
}

func replaceDict(src string, dict map[string]interface{}) string {
	for old, v := range dict {
		new, ok := v.(string)
		if !ok {
			continue
		}
		src = strings.Replace(src, old, new, -1)
	}
	return src
}
