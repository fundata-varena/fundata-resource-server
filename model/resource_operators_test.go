package model

import (
	"fmt"
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"testing"
)

func TestGetResourceUpdated(t *testing.T) {
	setup()
	_, err := GetResourceUpdated("", 213)
	fmt.Println(err)
}

func TestDownloadResource(t *testing.T) {
	setup()
	_ = DownloadResource("dota2_team_log", "hfhdfdjdhgh00")
}

func setup() {
	log.InitShareZapLogger(false)
	_ = conf.Init("../conf/dev.json")
}
