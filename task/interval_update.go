package task

import (
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"github.com/fundata-varena/fundata-resource-server/model"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

//
func IntervalUpdate() {
	config, err := conf.GetConf()
	if err != nil {

	}

	lock := semaphore.NewWeighted(1)

	ticker := time.NewTicker(time.Duration(config.Update.Interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			if !lock.TryAcquire(1) {
				continue
			}
			log.ShareZapLogger().Debug("IntervalUpdate start")
			process()
			log.ShareZapLogger().Debug("IntervalUpdate done")
			lock.Release(1)
		}
	}
}

//
func process() {
	ops := new(model.ResourceOps)

	// 最近一段时间内的更新
	rows, err := ops.GetResourceUpdated("", 123)
	if err != nil {
		return
	}

	if len(rows) == 0 {
		log.ShareZapLogger().Info("any resource updated")
		return
	}

	// 同步到本地
	var wg sync.WaitGroup

	for _, row := range rows {
		wg.Add(1)
		go func(r *model.ResourceUpdated) {
			defer wg.Done()
			log.ShareZapLogger().Debug(
				"Downloading",
				zap.String("resource_type", r.ResourceType),
				zap.String("resource_id", r.ResourceID))
			// 服务端的更新时间记录在本地
			err := ops.DownloadResource(r.ResourceType, r.ResourceID, r.UpdatedTime)
			if err != nil {
				log.ShareZapLogger().Error("DownloadResource err", zap.Error(err))
			}
		}(row)
	}

	wg.Wait()

	log.ShareZapLogger().Info("updated resources sync done")

	return
}