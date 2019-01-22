package fmap

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strings"
)

func goTypeWithGoPackage(p *descriptor.FileDescriptorProto, f *descriptor.FieldDescriptorProto) string {
	pkg := ""
	if *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE || *f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM {
		if isTimestampPackage(*f.TypeName) {
			pkg = "timestamp"
		} else {
			pkg = *p.GetOptions().GoPackage
			if strings.Contains(*p.GetOptions().GoPackage, ";") {
				pkg = strings.Split(*p.GetOptions().GoPackage, ";")[1]
			}
		}
	}
	return goTypeWithEmbedded(pkg, f, p)
}

// Warning does not handle message embedded like goTypeWithGoPackage does.
func goTypeWithPackage(f *descriptor.FieldDescriptorProto) string {
	pkg := ""
	if *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE || *f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM {
		if isTimestampPackage(*f.TypeName) {
			pkg = "timestamp"
		} else {
			pkg = getPackageTypeName(*f.TypeName)
		}
	}
	return goType(pkg, f)
}

func haskellType(pkg string, f *descriptor.FieldDescriptorProto) string {
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
			return fmt.Sprintf("[%s%s]", pkg, shortType(*f.TypeName))
		}
		return fmt.Sprintf("%s%s", pkg, shortType(*f.TypeName))
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[Word8]"
		}
		return "Word8"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return fmt.Sprintf("%s%s", pkg, shortType(*f.TypeName))
	default:
		return "Generic"
	}
}

func goTypeWithEmbedded(pkg string, f *descriptor.FieldDescriptorProto, p *descriptor.FileDescriptorProto) string {
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

			return fmt.Sprintf("[]*%s%s", pkg, shortType(name))
		}
		return fmt.Sprintf("*%s%s", pkg, shortType(name))
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
		return fmt.Sprintf("*%s%s", pkg, shortType(name))
	default:
		return "interface{}"
	}
}

//Deprecated. Instead use goTypeWithEmbedded
func goType(pkg string, f *descriptor.FieldDescriptorProto) string {
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
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return fmt.Sprintf("[]*%s%s", pkg, shortType(*f.TypeName))
		}
		return fmt.Sprintf("*%s%s", pkg, shortType(*f.TypeName))
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		if *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]byte"
		}
		return "byte"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return fmt.Sprintf("*%s%s", pkg, shortType(*f.TypeName))
	default:
		return "interface{}"
	}
}

func goZeroValue(f *descriptor.FieldDescriptorProto) string {
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

func jsType(f *descriptor.FieldDescriptorProto) string {
	template := "%s"
	if isFieldRepeated(f) {
		template = "Array<%s>"
	}

	switch *f.Type {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_ENUM:
		return fmt.Sprintf(template, namespacedFlowType(*f.TypeName))
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

func jsSuffixReservedKeyword(s string) string {
	return jsReservedRe.ReplaceAllString(s, "${1}${2}_${3}")
}
