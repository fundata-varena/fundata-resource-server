package model

import (
	"fmt"
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"testing"
	"time"
)

func TestGetResourceUpdated(t *testing.T) {
	setup()
	ops := new(ResourceOps)
	_, err := ops.GetResourceUpdated("", 213)
	fmt.Println(err)
}

func TestDownloadResource(t *testing.T) {
	setup()
	ops := new(ResourceOps)
	tm, _ := time.Parse("2006-01-02 15:04:05", "20190708120000")
	_ = ops.DownloadResource("dota2_team_log", "hfhdfdjdhgh00", tm)
}

func setup() {
	log.InitShareZapLogger(false)
	_ = conf.Init("../conf/dev.json")
}
