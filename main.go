package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/yuansudong/gengo"
)

// PluginName 用于描述一个插件名称
const PluginName = "protoc-gen-go-code:"

var (
	gomod string
)

func main() {
	flag.Parse()
	// SaveReq()
	startFromStdin()
}

// startFromStdin 用于开始从标准输入
func startFromStdin() {
	in := os.Stdin
	req, err := gengo.GetRequest(in)
	if err != nil {
		gengo.WriteError(err)
		return
	}
	if req.Parameter != nil {
		sParam := req.GetParameter()
		aParam := strings.Split(sParam, ";")
		if len(aParam) == 0 {
			log.Fatalln("need go_mod arguments")
		}
		for _, sPair := range aParam {
			aPair := strings.Split(sPair, "=")
			if len(aPair) != 2 {
				log.Fatalln("this is error synatx ", sPair)
			}
			sKey, sVal := aPair[0], aPair[1]
			switch sKey {
			case "go_mod":
				gomod = sVal
			}
		}

	}
	reg := gengo.NewRegistry()
	if err = reg.Load(req); err != nil {
		gengo.WriteError(err)
		return
	}
	var targets []*gengo.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			glog.Fatal(err)
		}
		targets = append(targets, f)
	}
	g := New(reg)
	out, err := g.Generate(targets)
	if err != nil {
		gengo.WriteError(err)
		return
	}
	gengo.WriteFiles(out)
}
