package main

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	link "gitlab.com/pufs/hugo-langs-symlink"
	"gitlab.com/pufs/hugo-langs-symlink/symlinkbuilder"
)

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter())
	builder := symlinkbuilder.New(logger.With().Str("compoent", "builder").Logger())

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create symlinks",
				Action: func(c *cli.Context) error {
					engine, err := newEngine(c, logger, builder)
					if err != nil {
						return err
					}
					if err := engine.Create(); err != nil {
						return err
					}
					logger.Info().Msg("symlinks are created")
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove symlinks",
				Action: func(c *cli.Context) error {
					engine, err := newEngine(c, logger, builder)
					if err != nil {
						return err
					}
					if err := engine.Remove(); err != nil {
						return err
					}
					logger.Info().Msg("symlinks are removed")
					return nil
				},
			},
		},
		Name:        "hslink",
		Description: "A tool to create symlinks for hugo multi-langs sites",
		Usage:       "A tool to create symlinks for hugo multi-langs sites",
		Authors: []*cli.Author{
			{
				Name:  "Cherie Hsieh",
				Email: "cherie@pufsecurity.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "langs",
				Value: nil,
				Usage: "available languages of the site",
			},
			&cli.StringSliceFlag{
				Name:  "skips",
				Value: nil,
				Usage: "the dirs you want to skip",
			},
			&cli.StringFlag{
				Name:     "dir",
				Value:    "",
				Usage:    "the content dir",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "ext",
				Value: ".md",
				Usage: "the file extension of content files. Default is `.md`",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newEngine(
	c *cli.Context,
	logger zerolog.Logger,
	builder link.SymlinkBuilder,
) (*link.Engine, error) {
	var opts []link.Option
	langs := c.StringSlice("langs")
	if langs != nil {
		opts = append(opts, link.SetLangs(langs))
	}
	skips := c.StringSlice("skips")
	if langs != nil {
		opts = append(opts, link.SetSkipDir(skips))
	}
	dir := c.String("dir")
	opts = append(opts, link.SetContentDir(dir))

	ext := c.String("ext")
	opts = append(opts, link.SetTargetExt(ext))

	cfg, err := link.NewConfiguration(opts...)
	if err != nil {
		return nil, err
	}

	logger = logger.With().
		Str("dir", dir).
		Strs("langs", langs).
		Strs("skips", skips).
		Str("extension", ext).Logger()
	return link.New(logger, builder, cfg), nil
}
