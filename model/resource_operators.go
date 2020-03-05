package model

import (
	"errors"
	"fmt"
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-go-sdk/fundata"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"github.com/fundata-varena/fundata-resource-server/database/mysql"
	"github.com/fundata-varena/fundata-resource-server/storage"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type ResourceOps struct {

}

// 本地单个资源获取
func (ops *ResourceOps) GetResource(resourceType, id string) (*ResourceLocal, error) {
	db, err := mysql.GetInstance()
	if err != nil {
		return nil, err
	}

	row := ResourceLocal{}
	has, err := db.Where("resource_type=? AND identifier=?", resourceType, id).Get(&row)
	if !has {
		return nil, nil
	}

	return &row, nil
}

func (ops *ResourceOps) GetLastUpdateTime() (*ResourceLocal, error) {
	db, err := mysql.GetInstance()
	if err != nil {
		return nil, err
	}

	row := ResourceLocal{}
	has, err := db.Where("id > 0").Limit(1).Desc("update_time").Get(&row)
	if !has {
		return nil, nil
	}

	return &row, nil
}

// 数据入库
// 更新时间存的是服务端的更新时间，并非生成记录的时间
func (ops *ResourceOps) InsertOrUpdate(resourceType, resourceId, savedAt string, updateTime time.Time) error {
	db, err := mysql.GetInstance()
	if err != nil {
		return err
	}

	row := ResourceLocal{}
	has, err := db.Where("resource_type=? AND identifier=?", resourceType, resourceId).Get(&row)
	if has {
		row.DstPath = savedAt
		row.UpdateTime = updateTime
		_, err = db.Id(row.Id).Update(&row)
	} else {
		row.ResourceType = resourceType
		row.Identifier = resourceId
		row.DstPath = savedAt
		row.AddTime = time.Now()
		row.UpdateTime = updateTime
		_, err = db.InsertOne(row)
	}
	if err != nil {
		return err
	}
	return nil
}

// 获取资源更新列表，for task
// after Unix时间戳，默认为-1
func (ops *ResourceOps) GetResourceUpdated(resourceType string, after int64, page, pageSize int) ([]*ResourceUpdated, error) {
	config, err := conf.GetConf()
	if err != nil {
		return nil, err
	}

	fundata.InitClient(config.ResourceService.Key, config.ResourceService.Secret)
	params := map[string]interface{}{}
	if resourceType != "" {
		params["resource_type"] = resourceType
	}
	if after > -1 {
		params["after"] = strconv.FormatInt(after,10)
	}
	params["page"] = page
	params["page_size"] = pageSize

	log.ShareZapLogger().Info("GetResourceUpdated request")
	resp, err := fundata.Get(config.ResourceService.UpdateListURI, params)
	log.ShareZapLogger().Info("GetResourceUpdated response", zap.Any("code", resp.RetCode))
	if err != nil {
		return nil, err
	}

	rows, ok := resp.Data.([]interface{})
	if !ok {
		// 不符合要求的数据即不需要处理的数据
		return nil, errors.New("empty result")
	}

	var results []*ResourceUpdated

	for _, value := range rows {
		row, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		res := new(ResourceUpdated)
		badRow := false
		for key, value := range row {
			vStr, ok := value.(string)
			if !ok {
				badRow = true
				break
			}
			if key == "resource_type" {
				res.ResourceType = vStr
			}
			if key == "resource_id" {
				res.ResourceID = vStr
			}
			if key == "size" {
				res.Size = vStr
			}
			if key == "updated_time" {
				tmInt64, err := strconv.ParseInt(vStr, 10, 64)
				if err != nil {
					log.ShareZapLogger().Error("time.Parse err", zap.Error(err))
					badRow = true
					break
				}
				res.UpdatedTime = time.Unix(tmInt64, 0)
			}
		}
		if !badRow {
			results = append(results, res)
		}
	}
	
	return results, nil
}

// 获取下载链接并转存到本地，for task
func (ops *ResourceOps) DownloadResource(resourceType, resourceId string, updateTime time.Time) error {
	config, err := conf.GetConf()
	if err != nil {
		return err
	}

	if resourceType == "" || resourceId == "" {
		return errors.New("resourceType & resourceId required")
	}

	// 拿到地址
	fundata.InitClient(config.ResourceService.Key, config.ResourceService.Secret)
	params := map[string]interface{}{}
	if resourceType != "" {
		params["resource_type"] = resourceType
	}
	if resourceId != "" {
		params["resource_id"] = resourceId
	}
	log.ShareZapLogger().Info("DownloadResource request start")
	resp, err := fundata.Get(config.ResourceService.DownloadURI, params)
	log.ShareZapLogger().Info("DownloadResource response", zap.Any("response", resp))
	if err != nil {
		return errors.New("")
	}

	resource, ok := resp.Data.(map[string]interface{})
	if !ok {
		return errors.New("response data incorrect")
	}

	valueInterface, ok := resource["url"]
	if !ok {
		return errors.New("resource url not set")
	}

	url, ok := valueInterface.(string)
	if !ok {
		return errors.New("resource url not string")
	}

	// 存放目录
	path := fmt.Sprintf("/%s", resourceType)
	// 组装本地文件名
	dstFileName := fmt.Sprintf("%s", resourceId)
	// 下载&转存
	savedAt, err := downloadAndSave(url, path, dstFileName)
	if err != nil {
		log.ShareZapLogger().Error(
			"Download and save err",
			zap.Any("error", err),
			zap.String("resource_url", url))
		return errors.New("downloadAndSave resource err")
	}

	log.ShareZapLogger().Debug("saved at ", zap.String("@", savedAt))

	// 写数据库
	err = ops.InsertOrUpdate(resourceType, resourceId, savedAt, updateTime)
	if err != nil {
		return err
	}

	time.Sleep(200 * time.Millisecond)

	return nil
}

func downloadAndSave(src, path, dstFileName string) (string, error) {
	// 下载
	client := http.Client{Timeout: 900 * time.Second}
	resp, err := client.Get(src)
	defer func() {
		_ = resp.Body.Close()
	}()

	// 如果是svg的，需要保存后缀
	if values, ok := resp.Header["Content-Type"]; ok && values[0] == "image/svg+xml" {
		dstFileName = dstFileName + ".svg"
	}

	// 转存
	storageIns, err := storage.GetInstance()
	if err != nil {
		return "", err
	}
	savedAt, err := storageIns.Store(resp.Body, path, dstFileName)
	if err != nil {
		return "", err
	}

	return savedAt, nil
}