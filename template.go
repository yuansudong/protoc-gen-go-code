package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/yuansudong/gengo"
	"github.com/yuansudong/protoc-gen-go-code/pboptions"
	"google.golang.org/protobuf/proto"

	plugin "github.com/yuansudong/gengo/plugin"
)

const codeFileTemplate = `
// 以此等代码注释开头的,都属于工具生成.请不要人为改变.
// 寄语:  工作是为了生活,而不是生活为了工作. 如何能在最少的时间内完成工作,这就是工具存在的意义!
///////////////////////////////////////////////////////////////////////////
// 		   Package: {{ .PackageName }}
// 		   Description: 服务于掌沃无限时制作
//		   Author: yuansudong
// 		   Protoc: unknown
//		   UpdateTime: {{ .DateTime }}
// 		   Company: 掌沃无限
//		   CreateTime: 2021	
///////////////////////////////////////////////////////////////////////////
package {{ .PackageName }}
{{ if .IsHave }}
import (
	"context"
	"sync"
	"github.com/yuansudong/gmiddleware"
	"{{.Gomod}}/core/state"
	"{{.Gomod}}/protocol/generate/gopb/error/pbm_bad"
	"{{.Gomod}}/protocol/generate/gopb/error/pbm_unimpl"
	{{ range $importKey,$importVal :=  .Imports  }}
		"{{ $importKey }}"
	{{ else }}
	{{ end }}
)
{{ else }}
{{ end }}

{{ range $serviceKey,$serviceVal :=  .Services }}
// 全局变量定义处
var xxx_{{ $serviceVal.Name }}_once sync.Once

var xxx_{{ $serviceVal.Name }}_inst *Impl{{ $serviceVal.Name }}Service

{{ range $methodKey,$methodVal := $serviceVal.Methods }}
	{{ if eq $methodVal.Mode 0 }}
	   type {{$serviceVal.Name }}{{ $methodVal.Name }}UnaryFn func(context.Context, *{{ $methodVal.TReq }})(*{{ $methodVal.TRsp }},error)
    {{ else if eq $methodVal.Mode 1}}
	    type {{$serviceVal.Name }}{{$methodVal.Name}}CsFn func({{$serviceVal.Name }}_{{ $methodVal.Name}}Server) error 
	{{ else }}
        type {{$serviceVal.Name }}{{ $methodVal.Name}}SsFn func(*{{ $methodVal.TReq }} , {{$serviceVal.Name }}_{{ $methodVal.Name}}Server) error
	{{end}}
{{ else }}
{{ end  }}


type Impl{{ $serviceVal.Name }}Service struct {
	chain           gmiddleware.ChanFn
{{ range $methodKey,$methodVal := $serviceVal.Methods }}
	{{if eq $methodVal.Mode 0}}
	    // Unary
		Real{{ $serviceVal.Name }}{{ $methodVal.Name }}UnaryFn {{ $serviceVal.Name }}{{ $methodVal.Name }}UnaryFn
		Info{{ $serviceVal.Name }}{{ $methodVal.Name }}    *gmiddleware.ServiceInfo
	{{else if eq $methodVal.Mode 1}}  
		// Client Stream
		Real{{ $serviceVal.Name }}{{ $methodVal.Name }}CsFn {{$serviceVal.Name }}{{$methodVal.Name}}CsFn
		Info{{ $serviceVal.Name }}{{ $methodVal.Name }}    *gmiddleware.ServiceInfo
	{{else}}
	    //  Service Stream
		Real{{ $serviceVal.Name }}{{ $methodVal.Name }}SsFn {{$serviceVal.Name }}{{$methodVal.Name}}SsFn
		Info{{ $serviceVal.Name }}{{ $methodVal.Name }}    *gmiddleware.ServiceInfo
	{{end}}
{{ else }}
{{ end  }}
}  

// NewImpl{{ $serviceVal.Name }}Service 实例化服务
func NewImpl{{ $serviceVal.Name }}Service() *Impl{{ $serviceVal.Name }}Service {
	inst := new(Impl{{ $serviceVal.Name }}Service)
{{ range $methodKey,$methodVal := $serviceVal.Methods }}
	{{ if eq $methodVal.Mode 0 }}
    inst.Real{{ $serviceVal.Name }}{{ $methodVal.Name }}UnaryFn = func(ctx context.Context, req *{{ $methodVal.TReq }})(rsp *{{ $methodVal.TRsp }},err error) {
		return nil, state.Unimplemented(pbm_unimpl.Unimpl_U_METHOD,"{{ $methodVal.Name }} Method Not Yet Implemented")
	}
	inst.Info{{ $serviceVal.Name }}{{ $methodVal.Name }} = &gmiddleware.ServiceInfo{ Service: "{{ $serviceVal.Name }}" ,Version: "{{ $serviceVal.Version }}", Method: "{{ $methodVal.Name }}"}
    {{ else if eq $methodVal.Mode 1}}
	inst.Real{{ $serviceVal.Name }}{{ $methodVal.Name }}CsFn = func(stream {{ $serviceVal.Name }}_{{ $methodVal.Name }}Server) error {
		return state.Unimplemented(pbm_unimpl.Unimpl_U_METHOD,"{{ $methodVal.Name }} Method Not Yet Implemented")
	}
	inst.Info{{ $serviceVal.Name }}{{ $methodVal.Name }} = &gmiddleware.ServiceInfo{ Service: "{{ $serviceVal.Name }}" ,Version: "{{ $serviceVal.Version }}", Method: "{{ $methodVal.Name }}"}
	{{ else }}
	inst.Real{{ $serviceVal.Name }}{{ $methodVal.Name }}SsFn = func(req *{{ $methodVal.TReq }},stream {{ $serviceVal.Name }}_{{ $methodVal.Name }}Server) error {
		return state.Unimplemented(pbm_unimpl.Unimpl_U_METHOD,"{{ $methodVal.Name }} Method Not Yet Implemented")
	}
	inst.Info{{ $serviceVal.Name }}{{ $methodVal.Name }} = &gmiddleware.ServiceInfo{ Service: "{{ $serviceVal.Name }}" ,Version: "{{ $serviceVal.Version }}", Method: "{{ $methodVal.Name }}"}
	{{end}}
{{ else }}
{{ end  }}
	return inst
} 
// WithChain 设置执行链
func (i *Impl{{ $serviceVal.Name }}Service) WithChain(chains ...gmiddleware.ChanFn) (*Impl{{ $serviceVal.Name }}Service) {
	i.chain = gmiddleware.ChainUnaryServer(chains...)
	return i
}
{{ range $methodKey,$methodVal := $serviceVal.Methods }}
{{ if eq $methodVal.Mode 0 }}
// unary模式 
func (i *Impl{{ $serviceVal.Name }}Service) With{{ $methodVal.Name }}(fn {{$serviceVal.Name }}{{ $methodVal.Name }}UnaryFn) *Impl{{ $serviceVal.Name }}Service {
	i.Real{{ $serviceVal.Name }}{{ $methodVal.Name }}UnaryFn = fn
	return i
}
// Deal{{ $methodVal.Name }} ....
func (i *Impl{{ $serviceVal.Name }}Service) Deal{{ $serviceVal.Name }}{{ $methodVal.Name }}(ctx context.Context, tmpReq interface{}) (interface{}, error) {
	req := tmpReq.(*{{ $methodVal.TReq }})
	if err := req.Validate(); err != nil {
		return nil, state.Bad(pbm_bad.Bad_BA_ARG, err.Error())
	}
	return i.Real{{$serviceVal.Name }}{{ $methodVal.Name }}UnaryFn(ctx, req)
}
// {{ $methodVal.Name }} 用于描述一个在线请求
func (i *Impl{{ $serviceVal.Name }}Service) {{ $methodVal.Name }}(ctx context.Context, req *{{ $methodVal.TReq }}) (rsp *{{ $methodVal.TRsp }}, err error) {
	tmpRsp, err := i.chain(ctx, req, i.Info{{ $serviceVal.Name }}{{ $methodVal.Name }}, i.Deal{{ $serviceVal.Name }}{{ $methodVal.Name }})
	if err != nil {
		return nil, err
	}
	return tmpRsp.(*{{ $methodVal.TRsp }}), nil
}
{{ else if eq $methodVal.Mode 1 }}

// 客户端流模式 
func (i *Impl{{ $serviceVal.Name }}Service) With{{ $methodVal.Name }}(fn {{$serviceVal.Name }}{{ $methodVal.Name }}CsFn) *Impl{{ $serviceVal.Name }}Service {
	i.Real{{ $serviceVal.Name }}{{ $methodVal.Name }}CsFn = fn
	return i
}

// Deal{{ $methodVal.Name }} ....
func (i *Impl{{ $serviceVal.Name }}Service) Deal{{ $serviceVal.Name }}{{ $methodVal.Name }}(stream {{$serviceVal.Name }}_{{ $methodVal.Name }}Server) error {
	return i.Real{{$serviceVal.Name }}{{ $methodVal.Name }}CsFn(stream)
}
// {{ $methodVal.Name }} 用于描述一个在线请求
func (i *Impl{{ $serviceVal.Name }}Service) {{ $methodVal.Name }}(stream {{$serviceVal.Name }}_{{ $methodVal.Name }}Server) error {
	err :=  i.Deal{{ $serviceVal.Name }}{{ $methodVal.Name }}(stream)
	if err != nil {
		return err
	}
	return nil
}
{{ else }}
// 服务端流模式 
func (i *Impl{{ $serviceVal.Name }}Service) With{{ $methodVal.Name }}(fn {{$serviceVal.Name }}{{ $methodVal.Name }}SsFn) *Impl{{ $serviceVal.Name }}Service {
	i.Real{{ $serviceVal.Name }}{{ $methodVal.Name }}SsFn = fn
	return i
}

// Deal{{ $methodVal.Name }} ....
func (i *Impl{{ $serviceVal.Name }}Service) Deal{{ $serviceVal.Name }}{{ $methodVal.Name }}(req *{{$methodVal.TReq }},stream {{$serviceVal.Name }}_{{ $methodVal.Name }}Server) error {
	return i.Real{{$serviceVal.Name }}{{ $methodVal.Name }}SsFn(req,stream)
}
// {{ $methodVal.Name }} 用于描述一个在线请求
func (i *Impl{{ $serviceVal.Name }}Service) {{$methodVal.Name}}(req *{{$methodVal.TReq }}, stream {{$serviceVal.Name }}_{{ $methodVal.Name }}Server) error {
	err :=  i.Deal{{ $serviceVal.Name }}{{ $methodVal.Name }}(req,stream)
	if err != nil {
		return err
	}
	return nil
}

{{ end }}


{{ else }}
{{ end }}

// GetImpl{{ $serviceVal.Name }} 用于获取一个实例
func GetImpl{{ $serviceVal.Name }}() *Impl{{ $serviceVal.Name }}Service {
	xxx_{{ $serviceVal.Name }}_once.Do(func(){
		xxx_{{ $serviceVal.Name }}_inst = NewImpl{{ $serviceVal.Name }}Service()
	})
	return xxx_{{ $serviceVal.Name }}_inst
}
{{ else }}
{{ end }}
`

// GoTemplate 用于生成相关的模板
func (g *gen) GoTemplate(file *gengo.File) *plugin.CodeGeneratorResponse_File {
	var err error
	as := &Args{}
	as.Imports = make(map[string]bool)
	as.PackageName = file.GoPkg.Name
	as.Gomod = gomod
	as.DateTime = time.Now().Local().String()
	buf := bytes.NewBuffer(make([]byte, 0, 40960))
	rspFile := new(plugin.CodeGeneratorResponse_File)
	for _, service := range file.Services {
		g.GoService(file, as, service)
	}
	if len(as.Services) != 0 {
		as.IsHave = true
	}

	tp := template.New("template.service")
	if tp, err = tp.Parse(codeFileTemplate); err != nil {
		log.Fatalln(err.Error())
	}
	if err = tp.Execute(buf, as); err != nil {
		log.Fatalln(err.Error())
	}
	name, err := g.GetAllFilePath(file)
	if err != nil {
		log.Println(PluginName, err.Error())
		os.Exit(-1)
	}
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	output := fmt.Sprintf("%s.pb.go-code.go", base)
	rspFile.Name = proto.String(output)
	rspFile.Content = proto.String(buf.String())
	return rspFile
}

//  GoMessage 用于处理消息
func (g *gen) GoService(rootFile *gengo.File, args *Args, service *gengo.Service) {
	sInst := Service{}

	if proto.HasExtension(rootFile.GetOptions(), pboptions.E_Openapiv2Swagger) {
		tag1 := proto.GetExtension(rootFile.GetOptions(), pboptions.E_Openapiv2Swagger).(*pboptions.Swagger)
		sInst.Version = tag1.Info.Version
	}
	sInst.Name = service.GetName()
	if proto.HasExtension(service.GetOptions(), pboptions.E_Openapiv2Tag) {
		tag1 := proto.GetExtension(service.GetOptions(), pboptions.E_Openapiv2Tag).(*pboptions.Tag)
		sInst.Name = tag1.Name
	}
	for _, mtd := range service.Methods {
		tMtd := Method{}
		if mtd.GetServerStreaming() {
			tMtd.Mode = ModeServiceStream
		}
		if mtd.GetClientStreaming() {
			tMtd.Mode = ModeClientStream
		}
		tMtd.Name = mtd.GetName()
		tMtd.TReq = mtd.RequestType.GoType(rootFile.GoPkg.Path)
		if mtd.RequestType.File.GoPkg.Path != rootFile.GoPkg.Path {
			args.Imports[mtd.RequestType.File.GoPkg.Path] = true
		}
		tMtd.TRsp = mtd.ResponseType.GoType(rootFile.GoPkg.Path)
		if mtd.ResponseType.File.GoPkg.Path != rootFile.GoPkg.Path {
			args.Imports[mtd.ResponseType.File.GoPkg.Path] = true
		}
		sInst.Methods = append(sInst.Methods, tMtd)
	}
	args.Services = append(args.Services, sInst)
}

func (g *gen) GetAllFilePath(file *gengo.File) (string, error) {
	name := file.GetName()
	switch {
	case g.modulePath != "" && g.pathType != pathTypeImport:
		return "", errors.New("cannot use module= with paths=")

	case g.modulePath != "":
		trimPath, pkgPath := g.modulePath+"/", file.GoPkg.Path+"/"
		if !strings.HasPrefix(pkgPath, trimPath) {
			return "", fmt.Errorf("%v: file go path does not match module prefix: %v", file.GoPkg.Path, trimPath)
		}
		return filepath.Join(strings.TrimPrefix(pkgPath, trimPath), filepath.Base(name)), nil

	case g.pathType == pathTypeImport && file.GoPkg.Path != "":
		return fmt.Sprintf("%s/%s", file.GoPkg.Path, filepath.Base(name)), nil

	default:
		return name, nil
	}
}
