package config

import (
	"fmt"
	"os"

	"github.com/supereagle/jenkins-pipeline/utils/json"
)

type Config struct {
	JenkinsServer       string `json:"jenkins_server,omitempty"`
	JenkinsUser         string `json:"jenkins_user,omitempty"`
	JenkinsPassword     string `json:"jenkins_password,omitempty"`
	JenkinsCredentialId string `json:"jenkins_credential,omitempty"`
}

func Read(path string) (*Config, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("The config file path should not be emtpy")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Fail to open the config file %s", path)
	}
	defer file.Close()

	cfg := &Config{}
	err = json.Unmarshal2JsonObj(file, cfg)
	if err != nil {
		return nil, fmt.Errorf("Fail to read the config file %s", path)
	}

	return cfg, nil
}
