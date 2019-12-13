package task

import (
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"github.com/fundata-varena/fundata-resource-server/model"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"time"
)

func IntervalUpdate(initialized bool) {
	config, err := conf.GetConf()
	if err != nil {
		log.ShareZapLogger().Error("IntervalUpdate get config nil")
		return
	}

	if !initialized {
		initData()
	}

	lock := semaphore.NewWeighted(1)

	ticker := time.NewTicker(time.Duration(config.Update.Interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			if !lock.TryAcquire(1) {
				continue
			}

			ops := new(model.ResourceOps)
			row, err := ops.GetLastUpdateTime()
			if err != nil {
				continue
			}
			if row == nil {
				continue
			}

			log.ShareZapLogger().Debug("IntervalUpdate start")
			process(row.UpdateTime.Unix())
			log.ShareZapLogger().Debug("IntervalUpdate done")

			lock.Release(1)
		}
	}
}

func initData() {
	// 初始化起始时间戳为1
	process(1)
}

func process(after int64) {
	ops := new(model.ResourceOps)

	page := 0
	pageSize := 20

	for true {
		log.ShareZapLogger().Info("processing", zap.Int("page", page))
		// 最近一段时间内的更新
		rows, err := ops.GetResourceUpdated("", after, page, pageSize)
		if err != nil {
			return
		}

		if len(rows) == 0 {
			break
		}

		for _, row := range rows {
			func(r *model.ResourceUpdated) {
				log.ShareZapLogger().Info("Downloading", zap.String("resource_type", r.ResourceType),
					zap.String("resource_id", r.ResourceID))
				// 服务端的更新时间记录在本地
				err := ops.DownloadResource(r.ResourceType, r.ResourceID, r.UpdatedTime)
				if err != nil {
					log.ShareZapLogger().Error("DownloadResource err", zap.Error(err))
				}
			}(row)
		}

		time.Sleep(200 * time.Millisecond)

		page++
	}

	log.ShareZapLogger().Info("updated resources sync done")

	return
}