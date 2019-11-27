package model

import (
	"time"
)

type ResourceLocal struct {
	Id           int64     `xorm:"pk autoincr BIGINT(20)"`
	ResourceType string    `xorm:"not null comment('资源类别') unique(res_id_uni_key) INT(11)"`
	Identifier    string    `xorm:"not null default '' comment('唯一标识') unique(res_id_uni_key) VARCHAR(255)"`
	DstPath      string    `xorm:"comment('目标OSS上的路径') VARCHAR(255)"`
	AddTime      time.Time `xorm:"default 'CURRENT_TIMESTAMP' comment('添加时间') TIMESTAMP"`
	UpdateTime   time.Time `xorm:"comment('更新时间') TIMESTAMP"`
}
