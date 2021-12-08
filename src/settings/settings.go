package settings

import (
	"elasticbulk/common"
	"fmt"
)

const APP_NAME = "elasticbulk"

type AppConfig struct {
	Elastic struct {
		Url string
		UserName string
		Password string
		Index string
		IndexName string
		Ip string
		Mask int
		NodeId string
		Offset int
		From string
		To string
		Territory string
	}
}

func(c AppConfig)String() string {
	return fmt.Sprintf("Appconfig-[Url:%v, UserName:%v, Password:%v]", c.Elastic.Url, c.Elastic.UserName, c.Elastic.Password)
}

var App AppConfig

func InitConfig() error {
	configDir := common.GetDefaultConfDir(APP_NAME)
	if err := common.ReadYamlConfig(configDir + "/" + APP_NAME + ".yml", &App); err != nil {
		return err
	}
	fmt.Printf("app config %v \n", App)
	return nil
}
