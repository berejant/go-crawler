package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const TestServerAddr = "localhost:44080"

func TestCrawler_Run(t *testing.T) {
	/**
	 * command to prepare `test-static` folder from `output` folder
	 * export DIR=ausland && mkdir -p ./test-static/${DIR} && cp output/www.spiegel.de_$(echo "$DIR" | sed 's/_/\//g').html ./test-static/${DIR}/index.html
	 *
	 * find test-static -name '*.html' -exec sed -i -e 's/https:\/\/www.spiegel.de\//http:\/\/localhost:44080\//g' {} \;
	 */

	t.Run("success", func(t *testing.T) {
		var expectedFilenames = []string{
			"localhost:44080.html",
			"localhost:44080_audio.html",
			"localhost:44080_fuermich.html",
			"localhost:44080_magazine.html",
			"localhost:44080_plus.html",
			"localhost:44080_politik_deutschland.html",
			"localhost:44080_schlagzeilen.html",
			"localhost:44080_spiegel.html",
			"localhost:44080_thema_klimawandel.html",
			"localhost:44080_thema_ukraine_konflikt.html",
		}

		var parser = &Parser{
			Fs: afero.NewMemMapFs(),
		}
		crawler.Parser = parser

		fs := http.FileServer(http.Dir("./test-static"))
		srv := &http.Server{
			Addr:    TestServerAddr,
			Handler: fs,
		}

		go func() {
			log.Print("Listening on " + srv.Addr + "...")
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()

		args := []string{
			"./crawler", "run",
			"--threads", "3",
			"--limit", strconv.Itoa(len(expectedFilenames)),
			"--verbose",
			"http://" + srv.Addr + "/",
		}

		fmt.Println("Run: " + strings.Join(args, " "))

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		App.Writer = stdout
		App.ErrWriter = stderr
		App.ExitErrHandler = func(cCtx *cli.Context, err error) {}

		err := App.Run(args)
		assert.NoError(t, err)

		if err = srv.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}

		files, _ := afero.ReadDir(parser.Fs, ".")
		assert.Len(t, files, len(expectedFilenames))

		actualFileNames := make([]string, 0, len(files))
		for _, file := range files {
			actualFileNames = append(actualFileNames, file.Name())
		}

		assert.Equal(t, expectedFilenames, actualFileNames)

		outString := stdout.String()
		assert.Contains(t, outString, "Ignore outside url ")
		assert.Contains(t, outString, "Reach url limit. Ignore url ")
	})

	t.Run("stop-on-too-many-errors", func(t *testing.T) {
		var parser = &Parser{
			Fs: afero.NewMemMapFs(),
		}
		crawler.Parser = parser

		fs := http.FileServer(http.Dir("./test-static"))
		srv := &http.Server{
			Addr:    TestServerAddr,
			Handler: fs,
		}

		go func() {
			log.Print("Listening on " + srv.Addr + "...")
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()

		args := []string{
			"./crawler", "run",
			"--threads", "3",
			"--limit", "50",
			"--verbose",
			"http://" + srv.Addr + "/404/",
		}

		fmt.Println("Run: " + strings.Join(args, " "))

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		App.Writer = stdout
		App.ErrWriter = stderr
		App.ExitErrHandler = func(cCtx *cli.Context, err error) {}

		err := App.Run(args)

		assert.Error(t, err)
		assert.Equal(t, ErrorTooManyError.Error(), err.Error())

		if err = srv.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}

		files, _ := afero.ReadDir(parser.Fs, ".")
		assert.Len(t, files, 1)

		outString := stdout.String()
		assert.Contains(t, outString, "Failed to process url: http://localhost:44080/404/10.html; error: HTTP error 404 Not Found")
	})
}
