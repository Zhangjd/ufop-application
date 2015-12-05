package main

import (
	"fmt"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/log"
	"os"
	"ufop"
	"ufop/foo"
)

const (
	VERSION = "0.1"
)

func help() {
	fmt.Printf("Usage: qufop <UfopConfig>\n\nVERSION: %s\n\n", VERSION)
}

func setQiniuHosts() {
	conf.RS_HOST = "http://rs.qiniu.com"
}

func main() {
	log.SetOutput(os.Stdout)
	setQiniuHosts()

	args := os.Args
	argc := len(args)

	var configFilePath string

	switch argc {
	case 2:
		configFilePath = args[1]
	default:
		help()
		return
	}

	//load config
	ufopConf := &ufop.UfopConfig{}
	confErr := ufopConf.LoadFromFile(configFilePath)
	if confErr != nil {
		log.Error("load config file error,", confErr)
		return
	}

	ufopServ := ufop.NewServer(ufopConf)

	//register job handlers
	if err := ufopServ.RegisterJobHandler("foo.conf", &foo.VideoMerger{}); err != nil {
		log.Error(err)
	}

	//listen
	ufopServ.Listen()
}
