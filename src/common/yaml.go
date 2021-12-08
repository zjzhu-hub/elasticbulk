package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var yamlConfig map[string]interface{}

func init() {
	yamlConfig = make(map[string]interface{})
}

func GetDefaultConfDir(app string) string {
	if IsExist("/etc/" + app) {
		return "/etc/" + app
	}
	return "conf"
}

func InitYamlConfig(f string) error {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	tmpConfig := make(map[string]interface{})
	yaml.Unmarshal(data, &tmpConfig)
	for key, value :=range tmpConfig {
		yamlConfig[key] = value;
	}
	return nil
}

func ReadYamlConfig(f string, object interface{}) error {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, object)
	return err
}

func GetYamlString(key, defaultValue string) (result string) {
	if yamlConfig[key] != nil {
		var ok bool
		result, ok = yamlConfig[key].(string)
		if !ok {
			result = defaultValue
		}
	} else {
		result = defaultValue
	}
	return
}

func SetYamlString(key string, val string) {
	yamlConfig[key] = val
}

func GetYamlInt(key string, defaultValue int) (result int) {
	if yamlConfig[key] != nil {
		var ok bool
		result, ok = yamlConfig[key].(int)
		if !ok {
			result = defaultValue
		}
	} else {
		result = defaultValue
	}
	return
}

func SetYamlInt(key string, val int) {
	yamlConfig[key] = val
}
