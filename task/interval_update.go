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
func IntervalUpdate(config *conf.Conf) {

	lock := semaphore.NewWeighted(1)

	ticker := time.NewTicker(time.Duration(config.Update.Interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			if !lock.TryAcquire(1) {
				continue
			}
			log.ShareZapLogger().Info("IntervalUpdate start")
			process()
			log.ShareZapLogger().Info("IntervalUpdate done")
			lock.Release(1)
		}
	}

}

//
func process() {
	// 最近一段时间内的更新
	rows, err := model.GetResourceUpdated("", 123)
	if err != nil {
		return
	}

	// 同步到本地
	var wg sync.WaitGroup

	for _, row := range rows {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.ShareZapLogger().Debug(
				"Downloading",
				zap.String("resource_type", row.ResourceType),
				zap.String("resource_id", row.ResourceID))
			//model.DownloadResource("", "")
		}()
	}

	wg.Wait()

	return
}