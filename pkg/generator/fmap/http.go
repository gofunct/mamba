package fmap

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	ggdescriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"strings"
)

func httpPath(m *descriptor.MethodDescriptorProto) string {

	ext, err := proto.GetExtension(m.Options, options.E_Http)
	if err != nil {
		return err.Error()
	}
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		return fmt.Sprintf("extension is %T; want an HttpRule", ext)
	}

	switch t := opts.Pattern.(type) {
	default:
		return ""
	case *options.HttpRule_Get:
		return t.Get
	case *options.HttpRule_Post:
		return t.Post
	case *options.HttpRule_Put:
		return t.Put
	case *options.HttpRule_Delete:
		return t.Delete
	case *options.HttpRule_Patch:
		return t.Patch
	case *options.HttpRule_Custom:
		return t.Custom.Path
	}
}

func httpPathsAdditionalBindings(m *descriptor.MethodDescriptorProto) []string {
	ext, err := proto.GetExtension(m.Options, options.E_Http)
	if err != nil {
		panic(err.Error())
	}
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		panic(fmt.Sprintf("extension is %T; want an HttpRule", ext))
	}

	var httpPaths []string
	var optsAdditionalBindings = opts.GetAdditionalBindings()
	for _, optAdditionalBindings := range optsAdditionalBindings {
		switch t := optAdditionalBindings.Pattern.(type) {
		case *options.HttpRule_Get:
			httpPaths = append(httpPaths, t.Get)
		case *options.HttpRule_Post:
			httpPaths = append(httpPaths, t.Post)
		case *options.HttpRule_Put:
			httpPaths = append(httpPaths, t.Put)
		case *options.HttpRule_Delete:
			httpPaths = append(httpPaths, t.Delete)
		case *options.HttpRule_Patch:
			httpPaths = append(httpPaths, t.Patch)
		case *options.HttpRule_Custom:
			httpPaths = append(httpPaths, t.Custom.Path)
		default:
			// nothing
		}
	}

	return httpPaths
}

func httpVerb(m *descriptor.MethodDescriptorProto) string {

	ext, err := proto.GetExtension(m.Options, options.E_Http)
	if err != nil {
		return err.Error()
	}
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		return fmt.Sprintf("extension is %T; want an HttpRule", ext)
	}

	switch t := opts.Pattern.(type) {
	default:
		return ""
	case *options.HttpRule_Get:
		return "GET"
	case *options.HttpRule_Post:
		return "POST"
	case *options.HttpRule_Put:
		return "PUT"
	case *options.HttpRule_Delete:
		return "DELETE"
	case *options.HttpRule_Patch:
		return "PATCH"
	case *options.HttpRule_Custom:
		return t.Custom.Kind
	}
}

func httpBody(m *descriptor.MethodDescriptorProto) string {

	ext, err := proto.GetExtension(m.Options, options.E_Http)
	if err != nil {
		return err.Error()
	}
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		return fmt.Sprintf("extension is %T; want an HttpRule", ext)
	}
	return opts.Body
}

func urlHasVarsFromMessage(path string, d *ggdescriptor.Message) bool {
	for _, field := range d.Field {
		if !isFieldMessage(field) {
			if strings.Contains(path, fmt.Sprintf("{%s}", *field.Name)) {
				return true
			}
		}
	}

	return false
}
