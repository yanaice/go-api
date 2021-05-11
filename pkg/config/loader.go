package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

var Conf struct {
	Env   string `yaml:"env"`
	Debug bool   `yaml:"debug"`
	CORS  struct {
		AllowAll       bool     `yaml:"allow-all"`
		AllowedDomains []string `yaml:"allowed-domains"`
	} `yaml:"cors"`
	Docs struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"docs"`
	JWT struct {
		Key            string `yaml:"key"`
		SignKey        string `yaml:"sign-key"`
		Timeout        string `yaml:"timeout"`
		RefreshKey     string `yaml:"refresh-key"`
		RefreshTimeout string `yaml:"refresh-timeout"`
		APIKey         string `yaml:"api-key"`
	} `yaml:"jwt"`
	AuthTokenAPI struct {
		BaseURL string `yaml:"base-url"`
	} `yaml:"auth-token-api"`
}

func getConfigFilesToLoad(confFilesStr string) []string {
	confFiles := strings.Split(confFilesStr, ",")

	confFilesToLoad := make([]string, 0, len(confFiles))

	for _, confFile := range confFiles {
		confFile = strings.TrimSpace(confFile)
		if _, err := os.Stat(confFile); os.IsNotExist(err) {
			continue
		}
		confFilesToLoad = append(confFilesToLoad, confFile)
	}

	return confFilesToLoad
}

func configMergeFiles(filenames []string) ([]byte, error) {
	mergedConf := make(map[string]interface{})

	for _, filename := range filenames {
		confData, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		thisConf := make(map[string]interface{})

		if err := yaml.Unmarshal(confData, &thisConf); err != nil {
			return nil, err
		}

		if err := mergo.Merge(&mergedConf, thisConf, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	mergedConfBytes, err := yaml.Marshal(mergedConf)
	if err != nil {
		return nil, err
	}

	return mergedConfBytes, nil
}

func Init(confFilesList string, conf interface{}) error {
	if confFilesList == "" {
		// Default config path
		confFilesList = "./configs/config.master.yml, ./configs/config.local.yml, ./configs/config.secret.yml, ./configs/config.yml"
	}

	confFilesToLoad := getConfigFilesToLoad(confFilesList)
	mergedConf, err := configMergeFiles(confFilesToLoad)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(mergedConf, conf); err != nil {
		return err
	}

	if err := yaml.Unmarshal(mergedConf, &Conf); err != nil {
		return err
	}

	return nil
}
