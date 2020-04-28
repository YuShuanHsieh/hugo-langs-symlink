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
	var opts []link.Option
	logger := zerolog.New(zerolog.NewConsoleWriter())
	builder := symlinkbuilder.New(logger.With().Str("compoent", "builder").Logger())

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create symlinks",
				Action: func(c *cli.Context) error {
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

					cfg, err := link.NewConfiguration(opts...)
					if err != nil {
						return err
					}
					logger.Info().Msgf("%v", cfg)
					engine := link.New(logger, builder, cfg)
					return engine.Create()
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove symlinks",
				Action: func(c *cli.Context) error {
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

					cfg, err := link.NewConfiguration(opts...)
					if err != nil {
						return err
					}
					logger.Info().Msgf("%v", cfg)
					engine := link.New(logger, builder, cfg)
					return engine.Remove()
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
				Name:  "dir",
				Value: "",
				Usage: "the content dir",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
