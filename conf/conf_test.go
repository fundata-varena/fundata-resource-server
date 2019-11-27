package conf

import (
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConf(t *testing.T) {
	log.InitShareZapLogger(false)
	_, err := Init("./dev.json")
	assert.Equal(t, nil, err)
}
