package conf

import (
	"errors"
	vp_cofig "git.vpgame.cn/sh-team/vp-go-sponsors/config"
)

var share *Conf

type Conf struct {
	Mysql struct {
		Db       string `json:"db"`
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		User     string `json:"user"`
	} `json:"mysql"`
	ResourceService struct {
		DownloadURI   string `json:"download_uri"`
		Key           string `json:"key"`
		Secret        string `json:"secret"`
		UpdateListURI string `json:"update_list_uri"`
	} `json:"resource_service"`
	TestKey string `json:"test_key"`
	Update  struct {
		Interval int `json:"interval"`
	} `json:"update"`
}


func InitConf(confFile string) error {
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
