package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const defaultUrlLimit = 100
const defaultThreadCount = 5

const tooManyErrorThreshold = 10

var ErrorTooManyError = errors.New("crawler got too many error ")

// strip query params from URL
var removeQueryStringRegexp = regexp.MustCompile(`\?.*`)

type Crawler struct {
	UrlLimit     uint32
	ThreadsCount uint8

	Parser ParserInterface
	Out    io.Writer

	// set of discovered urls to avoid visit same pages
	discoveredUrls *Set
	urlQueue       *StringQueue

	baseUrl string

	processedCounter uint32

	// channel for dispatch url from urlQueue into Threads
	threadChannel chan string
	errorCount    uint32
}

func (crawler *Crawler) Run(mainCtx context.Context, baseUrl string) (err error) {
	_, err = url.ParseRequestURI(baseUrl)
	if err != nil {
		return err
	}

	if crawler.ThreadsCount <= 0 {
		crawler.ThreadsCount = defaultThreadCount
	}

	if crawler.UrlLimit <= 0 {
		crawler.UrlLimit = defaultUrlLimit
	}

	err = crawler.Parser.Init()
	if err != nil {
		return err
	}

	crawler.baseUrl = baseUrl
	if err != nil {
		return err
	}

	// memory limit 2 KB per each URL. Average one URL needs 200-600 bytes.
	crawler.discoveredUrls = NewSet(int(crawler.UrlLimit) * 2 * 1024)

	crawler.urlQueue = NewStringQueue(100)
	crawler.threadChannel = make(chan string)

	crawler.processedCounter = 0

	crawler.addUrl(crawler.baseUrl)

	ctx, cancelFunc := context.WithCancel(mainCtx)
	wg := &sync.WaitGroup{}

	wg.Add(int(crawler.ThreadsCount))
	for i := uint8(0); i < crawler.ThreadsCount; i++ {
		go crawler.runThread(ctx, wg)
	}

	err = crawler.dispatchUrlQueueToThreadChannel(ctx)
	cancelFunc()
	wg.Wait()

	close(crawler.threadChannel)
	crawler.threadChannel = nil

	return err
}

func (crawler *Crawler) dispatchUrlQueueToThreadChannel(ctx context.Context) error {
	var nextUrl string
	for crawler.processedCounter < crawler.urlQueue.Counter && ctx.Err() == nil {
		if crawler.errorCount >= crawler.urlQueue.Counter || crawler.errorCount >= tooManyErrorThreshold {
			return ErrorTooManyError
		}

		nextUrl = crawler.urlQueue.GetNext()
		if nextUrl != "" {
			select {
			case crawler.threadChannel <- nextUrl:
			case <-ctx.Done():
			}
		} else {
			// useful for single-core and when crawler reach finish (by limit or site contents)
			runtime.Gosched()
			time.Sleep(time.Millisecond * 50)
		}
	}

	return nil
}
func (crawler *Crawler) addUrl(url string) bool {
	// make URL canonical
	url = removeQueryStringRegexp.ReplaceAllString(url, "")

	urlBytes := []byte(url)
	if crawler.discoveredUrls.Add(urlBytes) {
		if !strings.HasPrefix(url, crawler.baseUrl) {
			crawler.println("Ignore outside url " + url)
			return false
		}

		if crawler.urlQueue.Counter >= crawler.UrlLimit {
			crawler.println("Reach url limit. Ignore url " + url)
			return false
		}

		crawler.urlQueue.Add(url)
		return true
	}
	return false
}

func (crawler *Crawler) runThread(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var urlToProcess string

	for {
		select {
		case urlToProcess = <-crawler.threadChannel:
			crawler.processUrl(urlToProcess)

		case <-ctx.Done():
			return
		}
	}
}

func (crawler *Crawler) processUrl(url string) {
	crawler.println("Handle url: " + url)
	parserResponse, err := crawler.Parser.ProcessPage(url)
	atomic.AddUint32(&crawler.processedCounter, 1)

	if err != nil {
		atomic.AddUint32(&crawler.errorCount, 1)
		crawler.println("Failed to process url: " + url + "; error: " + err.Error())

	} else {
		for _, newUrl := range parserResponse.Links {
			crawler.addUrl(newUrl)
		}

		if crawler.errorCount > 1 {
			// decrement error count
			atomic.AddUint32(&crawler.errorCount, ^uint32(0))
		}
	}

}

func (crawler *Crawler) println(string string) {
	if crawler.Out != nil {
		_, _ = fmt.Fprintln(crawler.Out, string)
	}
}
