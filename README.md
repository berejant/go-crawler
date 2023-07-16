[![Run Tests](https://github.com/berejant/go-crawler/actions/workflows/release.yaml/badge.svg)](https://github.com/berejant/go-crawler/actions/workflows/release.yaml)
[![codecov](https://codecov.io/gh/berejant/go-crawler/branch/main/graph/badge.svg?token=pt1A4XNjiC)](https://codecov.io/gh/berejant/go-crawler)

## Web crawler

### Features:
 - Multiple threads using Golang coroutines
 - No need for any external Event queue (such as Beanstalk) to communicate between processes. Uses RAM instead of TCP-network with Event queue daemon. 
 - Very fast. Save 5000 pages in 1-2 minutes (depending on CPU and network).
 - Remember already discovered URL with [Fastcache](https://github.com/VictoriaMetrics/fastcache) (key-value storage). Doesn't process the same page twice.
 - Ignore URL outside specific web-URL (web-node).
 - Use canonical URL (without search param after `?`)
 - Save HTML to folder `output`
 - Configurable with CLI flags `--threads` and `--limit`
 - Could be run in Google Cloud Functions and AWS Lambda (just need an update to in cloud context instead of CLI context)
 - Use Abstract FileSystem [Afero](https://pkg.go.dev/github.com/spf13/afero) - it allow to use Google Cloud Storage or AWS S3.

### To-do list
 - Accept relative URL in HTML links.
 - Add test cases to reach 95% code coverage.
 - Can't work in a distributed environment. Just one instance per website.
 - Filesystem I/O is the bottleneck. Need to implement some memory cache layer.

### Run using the image
```shell
## see version
docker run --rm ghcr.io/berejant/go-crawler:main --version
```
```shell
mkdir -p output && chmod 0777 output
## run crawler
docker run -v "$(pwd)/output:/output" --rm ghcr.io/berejant/go-crawler:main run --threads 50 --limit 1000 --verbose https://www.spiegel.de/
```

### Build and run with Docker

```shell
DOCKER_BUILDKIT=1 docker build -t crawler . 

## see version
docker run --rm crawler --version
## see help
docker run --rm crawler run --help

## run crawler
docker run -v "$(pwd)/output:/output" --rm crawler run --threads 50 --limit 100 --verbose https://www.spiegel.de/
```

### Build and run with Go
```shell
go mod download
go test  . -covermode atomic
go build -o crawler .

## see help
./crawler run --version
./crawler run --help

## run crawler
./crawler run --threads 100 --limit 100 --verbose https://www.spiegel.de/
```

### Unit test
```shell

go mod download
go test  . -covermode atomic
# Test coverage: 91.4% of statements
```

### Benchmark test
```shell
time ./crawler run --threads 50 --limit 5000 https://www.spiegel.de/
```

### Benchmark result
```
# Apple M2 (8 cpu)
43.14s user 5.55s system 103% cpu 46.832 total 
```

### Clear output dir
```shell
rm -rf output && mkdir -p output && chmod 0777 output && touch output/.gitkeep
```
