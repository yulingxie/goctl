package sng

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type YapiConfig struct {
	Token string `yaml:"token,omitempty"`
	Id    int    `yaml:"id,omitempty"`
}

type JenkinsConfig struct {
	User string `yaml:"user,omitempty"`
	Pass string `yaml:"pass,omitempty"`
}

type SqlmConfig struct {
	Node  string   `yaml:"node,omitempty"`
	Db    string   `yaml:"db,omitempty"`
	Table []string `yaml:"table,omitempty"`
}

type RedismConfig struct {
	Node  string   `yaml:"node,omitempty"`
	Table []string `yaml:"table,omitempty"`
}

type MongomConfig struct {
	Node string   `yaml:"node,omitempty"`
	Db   string   `yaml:"db,omitempty"`
	Coll []string `yaml:"coll,omitempty"`
}

type SngServiceProjectConfig struct {
	ServiceName    string         `yaml:"service_name,omitempty"`
	ServiceDir     string         `yaml:"service_dir,omitempty"`
	ServiceTestDir string         `yaml:"service_test_dir,omitempty"`
	GitlabToken    string         `yaml:"gitlab_token,omitempty"`
	Yapi           YapiConfig     `yaml:"yapi,omitempty"`
	Jenkins        JenkinsConfig  `yaml:"jenkins,omitempty"`
	Sql            []SqlmConfig   `yaml:"sql,omitempty"`
	Redis          []RedismConfig `yaml:"redis,omitempty"`
	Mongo          []MongomConfig `yaml:"mongo,omitempty"`
}

func LoadSngServiceProjectConfigFromFile(filename string) (*SngServiceProjectConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &SngServiceProjectConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
