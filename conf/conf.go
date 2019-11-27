package conf

import (
	"errors"
	vp_cofig "git.vpgame.cn/sh-team/vp-go-sponsors/config"
)

// 仅在程序启动时写
// 运行时只并发读，故不做锁保护处理
var share *Conf

type Conf struct {
	FileStorage struct {
		FilePath string `json:"file_path"`
	} `json:"file_storage"`
	Mysql struct {
		Db       string `json:"db"`
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		User     string `json:"user"`
	} `json:"mysql"`
	OssStorage      struct{} `json:"oss_storage"`
	ResourceService struct {
		DownloadURI   string `json:"download_uri"`
		Key           string `json:"key"`
		Secret        string `json:"secret"`
		UpdateListURI string `json:"update_list_uri"`
	} `json:"resource_service"`
	StorageUsing string `json:"storage_using"`
	TestKey      string `json:"test_key"`
	Update       struct {
		Interval int `json:"interval"`
	} `json:"update"`
}


func Init(confFile string) error {
	if confFile == "" {
		return errors.New("init config with empty confFile path")
	}
	var conf Conf
	err := vp_cofig.LoadJSON(confFile, &conf)
	if err != nil {
		return err
	}

	share = &conf

	return nil
}

func GetConf() (*Conf, error) {
	if share == nil {
		return share, errors.New("init config first please")
	}

	return share, nil
}
