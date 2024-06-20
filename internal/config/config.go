package config

import "github.com/spf13/viper"

type config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SslMode  string `yaml:"sslmode"`
		Port     int    `yaml:"port"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
}

func NewConfig(configType, path, filename string) (*config, error) {
	cfg := new(config)

	viper.AddConfigPath(path)
	viper.SetConfigFile(filename)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
