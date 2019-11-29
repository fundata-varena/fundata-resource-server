package main

import (
	"flag"
	"fmt"
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"github.com/fundata-varena/fundata-resource-server/database/mysql"
	"github.com/fundata-varena/fundata-resource-server/router"
	"github.com/fundata-varena/fundata-resource-server/storage"
	"github.com/fundata-varena/fundata-resource-server/task"
	"net/http"
	_ "net/http/pprof"
)

var (
	confFilePtr = flag.String("conf_file", "", "")
	updateAutoPtr = flag.Bool("update_auto", false, "")
)

// todo SDK http client host&port revert
// todo SDK http header revert

func main() {
	flag.Parse()

	go func() {
		_ = http.ListenAndServe("127.0.0.1:6060", nil)
	}()

	log.InitShareZapLogger(false)

	// 初始化配置
	err := conf.Init(*confFilePtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 根据配置实例化存储
	err = storage.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	// mysql实例化
	err = mysql.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	if *updateAutoPtr {
		// 根据配置执行更新任务
		go task.IntervalUpdate()
	}

	r := router.NewRouter()
	err = r.Engine.Run()
	if err != nil {
		return
	}
}
