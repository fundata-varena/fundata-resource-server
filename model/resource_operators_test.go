package model

import (
	"fmt"
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"github.com/fundata-varena/fundata-resource-server/database/mysql"
	"testing"
	"time"
)

func TestGetResourceUpdated(t *testing.T) {
	setup()
	ops := new(ResourceOps)
	_, err := ops.GetResourceUpdated("", 213, 0, 20)
	fmt.Println(err)
}

func TestResourceOps_GetLastUpdateTime(t *testing.T) {
	setup()
	ops := new(ResourceOps)
	row, err := ops.GetLastUpdateTime()
	fmt.Println(err)
	fmt.Println(row.UpdateTime.Unix())
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
	_ = mysql.Init()
}
