package config

import (
	"os"

	"github.com/sirupsen/logrus"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/yaml"
)

type Server struct {
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
}

type Config struct {
	Server    Server              `json:"server,omitempty" yaml:"server,omitempty"`
	Ingresses []netv1.IngressSpec `json:"ingresses,omitempty" yaml:"ingresses,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Server: Server{
			Port: 8080,
		},
	}
}

func (cfg *Config) Initial() *Config {

	b, _ := yaml.Marshal(cfg)

	f, _ := os.OpenFile("config.default.yml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer f.Close()

	f.Write(b)

	return cfg
}

func (cfg *Config) ReadConfig() {
	for _, f := range []string{"config.default.yml", "config.yml", "config.server.yml", "config.ing.yml"} {
		err := cfg.readconfig(f)
		if err != nil {
			logrus.Warnf("read file %s failed: %v", f, err)
			continue
		}
	}
}

func (cfg *Config) readconfig(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		logrus.Warnf("read config failed: %v", err)
		return err
	}

	return yaml.Unmarshal(b, cfg)
}
