package hslink

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

type SymlinkBuilder interface {
	Build(origin, path string) error
	Remove(path string) error
	ShouldRemove(path string, info os.FileInfo) bool
}

type Engine struct {
	logger  zerolog.Logger
	builder SymlinkBuilder
	config  Configuration
}

func New(logger zerolog.Logger, builder SymlinkBuilder, cfg Configuration) *Engine {
	return &Engine{
		logger:  logger,
		builder: builder,
		config:  cfg,
	}
}

// Create create symlinks
func (e *Engine) Create() error {
	e.logger.Info().Msg("Starting to create symlinks..")
	if err := os.Chdir(e.config.contentDir); err != nil {
		return err
	}
	m := make(map[string]struct{})
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			e.logger.Error().Str("path", path).Err(err)
			return err
		}
		if info.IsDir() {
			for _, v := range e.config.skipDirs {
				if info.Name() == v {
					return filepath.SkipDir
				}
			}
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			name := info.Name()
			if strings.HasSuffix(name, "_index.md") {
				return nil
			}
			if _, ok := m[name]; ok {
				return nil
			}
			name = getContentName(e.config.langs, name)
			if name == "" {
				return nil
			}

			for _, lang := range e.config.langs {
				fileName := fmt.Sprintf("%s.%s.md", name, lang)
				origin := filepath.Join(e.config.contentDir, path)
				expected := filepath.Join(e.config.contentDir, filepath.Dir(path), fileName)
				if _, err := os.Stat(expected); err != nil {
					err := e.builder.Build(origin, expected)
					if err != nil {
						e.logger.Error().Str("path", expected).Err(err)
					}
					m[fileName] = struct{}{}
				}
			}
		}
		return nil
	})
}

func (e *Engine) ShowConfig() {
	e.logger.Info().Msg("The configuration is:")
	e.logger.Info().Msgf("Dir - %s", e.config.contentDir)
	e.logger.Info().Msgf("Langs - %s", strings.Join(e.config.langs, ","))
	e.logger.Info().Msgf("Skips - %s", strings.Join(e.config.skipDirs, ","))
}

func (e *Engine) Remove() error {
	if err := os.Chdir(e.config.contentDir); err != nil {
		e.logger.Error().Str("path", e.config.contentDir).Err(err)
		return err
	}
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			e.logger.Error().Str("path", path).Err(err)
			return err
		}
		if info.IsDir() {
			for _, v := range e.config.skipDirs {
				if info.Name() == v {
					return filepath.SkipDir
				}
			}
		} else {
			if e.builder.ShouldRemove(path, info) {
				fPath := filepath.Join(e.config.contentDir, path)
				e.builder.Remove(fPath)
			}
		}
		return nil
	})
}

func getContentName(langs []string, name string) string {
	nameWithoutType := strings.TrimSuffix(name, ".md")
	for _, v := range langs {
		if strings.HasSuffix(nameWithoutType, v) {
			return strings.TrimSuffix(name, "."+v)
		}
	}
	return nameWithoutType
}
