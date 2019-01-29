package plugins

import (
	"bytes"
	"github.com/gofunct/mamba/plugins/fmap"
	"github.com/gofunct/mamba/plugins/fmap/protofunc"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type GenericTemplateBasedEncoder struct {
	templateDir    string
	service        *descriptor.ServiceDescriptorProto
	file           *descriptor.FileDescriptorProto
	enum           []*descriptor.EnumDescriptorProto
	debug          bool
	destinationDir string
}

type Ast struct {
	BuildDate      time.Time                          `json:"build-date"`
	BuildHostname  string                             `json:"build-hostname"`
	BuildUser      string                             `json:"build-user"`
	GoPWD          string                             `json:"go-pwd,omitempty"`
	PWD            string                             `json:"pwd"`
	Debug          bool                               `json:"debug"`
	DestinationDir string                             `json:"destination-dir"`
	File           *descriptor.FileDescriptorProto    `json:"file"`
	RawFilename    string                             `json:"raw-filename"`
	Filename       string                             `json:"filename"`
	TemplateDir    string                             `json:"template-dir"`
	Service        *descriptor.ServiceDescriptorProto `json:"service"`
	Enum           []*descriptor.EnumDescriptorProto  `json:"enum"`
}

func NewGenericServiceTemplateBasedEncoder(templateDir string, service *descriptor.ServiceDescriptorProto, file *descriptor.FileDescriptorProto, debug bool, destinationDir string) (e *GenericTemplateBasedEncoder) {
	e = &GenericTemplateBasedEncoder{
		service:        service,
		file:           file,
		templateDir:    templateDir,
		debug:          debug,
		destinationDir: destinationDir,
		enum:           file.GetEnumType(),
	}
	if debug {
		log.Printf("new encoder: file=%q service=%q template-dir=%q", file.GetName(), service.GetName(), templateDir)
	}
	InitPathMap(file)

	return
}

func NewGenericTemplateBasedEncoder(templateDir string, file *descriptor.FileDescriptorProto, debug bool, destinationDir string) (e *GenericTemplateBasedEncoder) {
	e = &GenericTemplateBasedEncoder{
		service:        nil,
		file:           file,
		templateDir:    templateDir,
		enum:           file.GetEnumType(),
		debug:          debug,
		destinationDir: destinationDir,
	}
	if debug {
		log.Printf("new encoder: file=%q template-dir=%q", file.GetName(), templateDir)
	}
	InitPathMap(file)

	return
}

func (e *GenericTemplateBasedEncoder) templates() ([]string, error) {
	filenames := []string{}

	err := filepath.Walk(e.templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}
		rel, err := filepath.Rel(e.templateDir, path)
		if err != nil {
			return err
		}
		if e.debug {
			log.Printf("new template: %q", rel)
		}

		filenames = append(filenames, rel)
		return nil
	})
	return filenames, err
}

func (e *GenericTemplateBasedEncoder) genAst(templateFilename string) (*Ast, error) {
	// prepare the ast passed to the template engine
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	goPwd := ""
	if os.Getenv("GOPATH") != "" {
		goPwd, err = filepath.Rel(os.Getenv("GOPATH")+"/src", pwd)
		if err != nil {
			return nil, err
		}
		if strings.Contains(goPwd, "../") {
			goPwd = ""
		}
	}
	ast := Ast{
		BuildDate:      time.Now(),
		BuildHostname:  hostname,
		BuildUser:      os.Getenv("USER"),
		PWD:            pwd,
		GoPWD:          goPwd,
		File:           e.file,
		TemplateDir:    e.templateDir,
		DestinationDir: e.destinationDir,
		RawFilename:    templateFilename,
		Filename:       "",
		Service:        e.service,
		Enum:           e.enum,
	}
	buffer := new(bytes.Buffer)

	unescaped, err := url.QueryUnescape(templateFilename)
	if err != nil {
		log.Printf("failed to unescape filepath %q: %v", templateFilename, err)
	} else {
		templateFilename = unescaped
	}

	tmpl, err := template.New("").Funcs(fmap.DefaultFMap).Parse(templateFilename)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buffer, ast); err != nil {
		return nil, err
	}
	ast.Filename = buffer.String()
	return &ast, nil
}

func (e *GenericTemplateBasedEncoder) buildContent(templateFilename string) (string, string, error) {
	// initialize template engine
	fullPath := filepath.Join(e.templateDir, templateFilename)
	templateName := filepath.Base(fullPath)
	tmpl, err := template.New(templateName).Funcs(fmap.DefaultFMap).ParseFiles(fullPath)
	if err != nil {
		return "", "", err
	}

	ast, err := e.genAst(templateFilename)
	if err != nil {
		return "", "", err
	}

	// generate the content
	buffer := new(bytes.Buffer)
	if err := tmpl.Execute(buffer, ast); err != nil {
		return "", "", err
	}

	return buffer.String(), ast.Filename, nil
}

func (e *GenericTemplateBasedEncoder) Files() []*plugin_go.CodeGeneratorResponse_File {
	templates, err := e.templates()
	if err != nil {
		log.Fatalf("cannot get templates from %q: %v", e.templateDir, err)
	}

	length := len(templates)
	files := make([]*plugin_go.CodeGeneratorResponse_File, 0, length)
	errChan := make(chan error, length)
	resultChan := make(chan *plugin_go.CodeGeneratorResponse_File, length)
	for _, templateFilename := range templates {
		go func(tmpl string) {
			var translatedFilename, content string
			content, translatedFilename, err = e.buildContent(tmpl)
			if err != nil {
				errChan <- err
				return
			}
			filename := translatedFilename[:len(translatedFilename)-len(".tmpl")]

			resultChan <- &plugin_go.CodeGeneratorResponse_File{
				Content: &content,
				Name:    &filename,
			}
		}(templateFilename)
	}
	for i := 0; i < length; i++ {
		select {
		case f := <-resultChan:
			files = append(files, f)
		case err = <-errChan:
		}
	}
	if err != nil {
		panic(err)
	}
	return files
}

func InitPathMap(file *descriptor.FileDescriptorProto) {
	protofunc.PathMap = make(map[interface{}]*descriptor.SourceCodeInfo_Location)
	addToPathMap(file.GetSourceCodeInfo(), file, []int32{})
}

func InitPathMaps(files []*descriptor.FileDescriptorProto) {
	protofunc.PathMap = make(map[interface{}]*descriptor.SourceCodeInfo_Location)
	for _, file := range files {
		addToPathMap(file.GetSourceCodeInfo(), file, []int32{})
	}
}

// addToPathMap traverses through the AST adding SourceCodeInfo_Location entries to the protofunc.PathMap.
// Since the AST is a tree, the recursion finishes once it has gone through all the nodes.
func addToPathMap(info *descriptor.SourceCodeInfo, i interface{}, path []int32) {
	loc := findLoc(info, path)
	if loc != nil {
		protofunc.PathMap[i] = loc
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
