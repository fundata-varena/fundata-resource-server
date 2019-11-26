package model

import (
	"errors"
	"fmt"
	"git.vpgame.cn/sh-team/vp-go-sponsors/log"
	"github.com/fundata-varena/fundata-go-sdk/fundata"
	"github.com/fundata-varena/fundata-resource-server/conf"
	"go.uber.org/zap"
	"strconv"
)

// 本地单个资源获取
func GetResource() {

}

// 本地批量获取
func GetResources() {

}

// 获取资源更新列表，for task
// after Unix时间戳，默认为-1
func GetResourceUpdated(resourceType string, after int64) ([]*ResourceUpdated, error) {
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
	log.ShareZapLogger().Debug("GetResourceUpdated request")
	resp, err := fundata.Get(config.ResourceService.UpdateListURI, params)
	log.ShareZapLogger().Debug("GetResourceUpdated response", zap.Any("response", resp))
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
				res.UpdatedTime = vStr
			}
		}
		if !badRow {
			results = append(results, res)
		}
	}
	
	return results, nil
}

// 获取下载链接并转存到本地，for task
func DownloadResource(resourceType, resourceId string) error {
	config, err := conf.GetConf()
	if err != nil {
		return err
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
	resp, err := fundata.Get(config.ResourceService.DownloadURI, params)

	fmt.Println(resp)

	// 下载

	// 转存

	return nil
}