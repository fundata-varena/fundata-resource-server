package main

import (
	"flag"
	"fmt"
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"github.com/fundata-varena/fundata-resource-server/database/mysql"
	"github.com/fundata-varena/fundata-resource-server/router"
	"github.com/fundata-varena/fundata-resource-server/task"
)

var (
	confFilePtr = flag.String("conf_file", "", "")
	updateAutoPtr = flag.Bool("update_auto", false, "")
)

func main() {
	flag.Parse()

	log.InitShareZapLogger(false)

	err := conf.InitConf(*confFilePtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	config, err := conf.GetConf()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = mysql.Init(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *updateAutoPtr {
		go task.IntervalUpdate(config)
	}

	r := router.NewRouter()
	err = r.Engine.Run()
	if err != nil {
		return
	}
}
