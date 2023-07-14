ARG GO_VERSION=${GO_VERSION:-1.19}

FROM --platform=${BUILDPLATFORM:-linux/amd64}  golang:${GO_VERSION}-alpine AS builder

WORKDIR /src/
RUN apk update && apk add --no-cache git
RUN cat /etc/passwd | grep nobody > /etc/passwd.nobody

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

# Build the binary.
RUN --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -tags=nomsgpack -o /crawler .

# build a small image
FROM --platform=${BUILDPLATFORM:-linux/amd64}  alpine

ENV TZ=Europe/Kyiv
RUN apk add tzdata

RUN mkdir /output && chmod 777 -R /output

COPY --from=builder /etc/passwd.nobody /etc/passwd
COPY --from=builder /crawler /crawler
WORKDIR /

# Run
USER nobody
ENTRYPOINT ["/crawler"]
# run --threads 10 --limit 200 https://www.spiegel.de/
