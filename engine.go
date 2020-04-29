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
	if err := os.Chdir(e.config.contentDir); err != nil {
		return err
	}
	m := make(map[string]struct{})
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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
		if !info.IsDir() && filepath.Ext(info.Name()) == e.config.targetExt {
			name := info.Name()
			if _, ok := m[name]; ok {
				return nil
			}
			name = e.parseFileName(name)
			if name == "" {
				return nil
			}

			for _, lang := range e.config.langs {
				fileName := fmt.Sprintf("%s.%s%s", name, lang, e.config.targetExt)
				origin := filepath.Join(e.config.contentDir, path)
				expected := filepath.Join(e.config.contentDir, filepath.Dir(path), fileName)
				if _, err := os.Stat(expected); err != nil {
					err := e.builder.Build(origin, expected)
					if err != nil {
						e.logger.Error().Str("path", expected).Err(err)
					} else {
						e.logger.Info().Msgf("symlink [%s] is created", expected)
					}
				}
				m[fileName] = struct{}{}
			}
		}
		return nil
	})
	return err
}

// Remove remove all symlinks
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
				if err := e.builder.Remove(fPath); err != nil {
					e.logger.Error().Err(err).Msgf("failed to remove the symlink [%s]", fPath)
				} else {
					e.logger.Info().Msgf("symlink [%s] is removed", fPath)
				}
			}
		}
		return nil
	})
}

func (e *Engine) parseFileName(name string) string {
	nameWithoutType := strings.TrimSuffix(name, e.config.targetExt)
	for _, v := range e.config.langs {
		if strings.HasSuffix(nameWithoutType, v) {
			return strings.TrimSuffix(nameWithoutType, "."+v)
		}
	}
	return nameWithoutType
}
