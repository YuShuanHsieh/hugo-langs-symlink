package hslink_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	hslink "github.com/YuShuanHsieh/hugo-langs-symlink"
	"github.com/YuShuanHsieh/hugo-langs-symlink/symlinkbuilder"
)

func TestHsLink(t *testing.T) {
	logger := zerolog.New(zerolog.NewConsoleWriter())
	builder := symlinkbuilder.New(logger)
	langs := []string{"uk", "tw"}
	filesName := []string{"about.md", "books.md", "books.tw.md"}

	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	defer os.RemoveAll(tmpDir)

	for _, name := range filesName {
		_, err := os.Create(filepath.Join(tmpDir, name))
		assert.NoError(t, err)
	}
	cfg, err := hslink.NewConfiguration(
		hslink.SetContentDir(tmpDir),
		hslink.SetLangs(langs),
		hslink.SetTargetExt(".md"),
	)
	assert.NoError(t, err)
	engine := hslink.New(logger, builder, cfg)
	engine.Create()

	assert.True(t, checkFilesExsistence(tmpDir, "about.uk.md"))
	assert.True(t, checkFilesExsistence(tmpDir, "about.tw.md"))
	assert.True(t, checkFilesExsistence(tmpDir, "books.uk.md"))

	engine.Remove()
	assert.False(t, checkFilesExsistence(tmpDir, "about.uk.md"))
	assert.False(t, checkFilesExsistence(tmpDir, "about.tw.md"))
	assert.False(t, checkFilesExsistence(tmpDir, "books.uk.md"))
}

func checkFilesExsistence(tmpDir, file string) bool {
	fPath := filepath.Join(tmpDir, file)
	if _, err := os.Stat(fPath); err != nil {
		return false
	}
	return true
}
