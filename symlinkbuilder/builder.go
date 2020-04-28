package symlinkbuilder

import (
	"os"

	"github.com/rs/zerolog"
)

type Symlinkbuilder struct {
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *Symlinkbuilder {
	return &Symlinkbuilder{
		logger: logger,
	}
}

func (t *Symlinkbuilder) Build(origin, path string) error {
	return os.Symlink(origin, path)
}

func (t *Symlinkbuilder) Remove(path string) error {
	return os.Remove(path)
}

func (t *Symlinkbuilder) ShouldRemove(path string, info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}
