package main

import "context"

type CrawlerInterface interface {
	Run(ctx context.Context, startUrl string) error
}
