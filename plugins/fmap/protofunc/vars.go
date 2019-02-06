package protofunc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	ggdescriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"regexp"
)

var (
	PathMap      map[interface{}]*descriptor.SourceCodeInfo_Location
	registry     *ggdescriptor.Registry // some helpers need access to registry
	jsReservedRe = regexp.MustCompile(`(^|[^A-Za-z])(do|if|in|for|let|new|try|var|case|else|enum|eval|false|null|this|true|void|with|break|catch|class|const|super|throw|while|yield|delete|export|import|public|return|static|switch|typeof|default|extends|finally|package|private|continue|debugger|function|arguments|interface|protected|implements|instanceof)($|[^A-Za-z])`)
)
