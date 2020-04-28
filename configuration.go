package hslink

import (
	"errors"
	"os"
)

// Configuration config params
type Configuration struct {
	langs      []string
	skipDirs   []string
	contentDir string
}

type Option func(cfg *Configuration) error

func NewConfiguration(opts ...Option) (Configuration, error) {
	var cfg Configuration
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func SetContentDir(path string) Option {
	return func(cfg *Configuration) error {
		if info, err := os.Stat(path); err != nil {
			return err
		} else if !info.IsDir() {
			return errors.New("the path should be a dir")
		}
		cfg.contentDir = path
		return nil
	}
}

func SetLangs(langs []string) Option {
	return func(cfg *Configuration) error {
		cfg.langs = langs
		return nil
	}
}

func SetSkipDir(skipDirs []string) Option {
	return func(cfg *Configuration) error {
		cfg.skipDirs = skipDirs
		return nil
	}
}
