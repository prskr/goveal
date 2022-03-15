package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Load(workingDir, configFile string) (cfg *Components, err error) {
	var (
		loader = viper.New()
		home   string
	)
	cfg = new(Components)

	for k, v := range defaults {
		loader.SetDefault(k, v)
	}

	if configFile != "" {
		loader.SetConfigFile(configFile)
	} else if home, err = homedir.Dir(); err != nil {
		return nil, err
	} else {
		loader.AddConfigPath(home)
		loader.AddConfigPath(workingDir)
		loader.SetConfigName("goveal")
		loader.SetConfigType("yaml")
	}

	loader.AutomaticEnv()

	if err = loader.ReadInConfig(); err == nil {
		log.Info("Using config file:", loader.ConfigFileUsed())
		cfg.ConfigFileInUse = loader.ConfigFileUsed()
		loader.WatchConfig()
	} else {
		return nil, err
	}

	loader.OnConfigChange(func(in fsnotify.Event) {
		if in.Op == fsnotify.Write {
			_ = loader.Unmarshal(cfg)
		}
	})

	err = loader.Unmarshal(cfg)

	return cfg, err
}
