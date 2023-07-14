package main

import (
	"context"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"os/signal"
	"syscall"
)

const defaultErrorExitCode = 1

var ErrCodeMap = map[error]int{
	ErrorTooManyError: 100,
}

var App = &cli.App{
	Name:           "crawler",
	Usage:          "webcrawler (by Anton Berezhnyi)",
	DefaultCommand: "run",
	Version:        "v0.1.0",
	Commands: []*cli.Command{
		{
			Name:      "run",
			Aliases:   []string{"r"},
			Usage:     "Run crawler",
			Action:    runCrawler,
			ArgsUsage: "URL",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:        "threads",
					Aliases:     []string{"T"},
					Usage:       "threads 5",
					DefaultText: "5",
				},

				&cli.IntFlag{
					Name:        "limit",
					Aliases:     []string{"L"},
					Usage:       "limit 100",
					DefaultText: "100",
				},

				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"V"},
					Usage:       "verbose",
					DefaultText: "false",
				},
			},
		},
	},
}

var crawler = &Crawler{
	Parser: &Parser{
		Fs: afero.NewBasePathFs(afero.NewOsFs(), "output"),
	},
}

func runCrawler(c *cli.Context) error {
	crawler.UrlLimit = uint32(c.Uint64("limit"))
	crawler.ThreadsCount = uint8(c.Uint64("threads"))

	if c.Bool("verbose") {
		crawler.Out = c.App.Writer
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	err := crawler.Run(ctx, c.Args().Get(0))

	if err != nil {
		errorCode := ErrCodeMap[err]
		if errorCode == 0 {
			errorCode = defaultErrorExitCode
		}

		return cli.Exit(err, errorCode)
	}
	return nil
}
