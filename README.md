[![Run Tests](https://github.com/berejant/go-crawler/actions/workflows/release.yaml/badge.svg)](https://github.com/berejant/go-crawler/actions/workflows/release.yaml)
[![codecov](https://codecov.io/gh/berejant/go-crawler/branch/main/graph/badge.svg?token=pt1A4XNjiC)](https://codecov.io/gh/berejant/go-crawler)

### Run using image
```shell
## see version
docker run --rm ghcr.io/berejant/go-crawler:main --version
```
```shell
## run crawler
docker run -v "$(pwd)/output:/output" --rm ghcr.io/berejant/go-crawler:main run --threads 50 --limit 1000 https://www.spiegel.de/
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


# Intel(R) Xeon(R) CPU E5-2609 v4 @ 1.70GHz (16 cpu)
 
```

### Clear output dir
```shell
rm -rf output && mkdir output && touch output/.gitkeep
```