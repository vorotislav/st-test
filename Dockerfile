FROM golang:1.22.1-alpine as builder

WORKDIR /src

COPY . /src

RUN apk --no-cache add bash git

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

ENV BUILD_VERSION="$(shell git describe --tags)"
ENV BUILD_DATE="$(shell date +%FT%T%z)"
ENV BUILD_COMMIT="$(shell git rev-parse --short HEAD)"
ENV PKG_PATH="./cmd/util"

ENV LDFLAGS="-X ${PKG_PATH}.buildVersion=$(BUILD_VERSION) -X ${PKG_PATH}.buildDate=$(BUILD_DATE) -X ${PKG_PATH}.buildCommit=$(BUILD_COMMIT)"

RUN CGO_ENABLED=0 go build -trimpath -ldflags "${LDFLAGS}" -o /tmp/st-test ./cmd

FROM alpine as runner

COPY --from=builder /tmp/st-test /
COPY ./bin/config.yaml /config.yaml

CMD ["./st-test", "config=config.yaml"]