package main

import (
	"github.com/yuansudong/gengo"
	plugin "github.com/yuansudong/gengo/plugin"
)

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)

type gen struct {
	reg        *gengo.Registry
	pathType   pathType
	modulePath string
}

type wrapper struct {
	fileName string
}

const (
	// ModeClientStream 客户端流模式
	ModeClientStream = 1
	// ModeClientStream 服务端流模式
	ModeServiceStream = 2
)

// Method 用于描述一个服务名称
type Method struct {
	Name string
	TReq string
	TRsp string
	Mode int
}

// Service 用于描述一个服务
type Service struct {
	Version string
	Methods []Method
	Name    string
}

// Args 用于描述一个参数
type Args struct {
	DateTime    string
	IsHave      bool
	PackageName string
	Services    []Service
	Imports     map[string]bool
	Gomod       string
}

// New returns a new generator which generates grpc gateway files.
func New(reg *gengo.Registry) *gen {
	return &gen{
		reg: reg,
	}
}

// Generate 用于执行生成操作
func (g *gen) Generate(targets []*gengo.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range targets {
		files = append(files, g.GoTemplate(file))
	}
	return files, nil
}
