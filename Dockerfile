FROM golang:1.16-alpine as builder

ARG CI_TAG
ARG BUILD_NUMBER
ARG BUILD_COMMIT_SHORT
ARG CI_BUILD_CREATED
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
ENV CGO_ENABLED=0
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/
RUN go build -o release/drone-kaniko \
    -ldflags "-w -s \
   -X main.version=${CI_TAG} \
   -X main.buildNumber=${BUILD_NUMBER} \
   -X main.commit=${BUILD_COMMIT_SHORT} \
   -X main.buildTime=${CI_BUILD_CREATED}" \
    -tags netgo \
    ./cmd/drone-kaniko

# use the debug image since it comes with /kaniko/warmer
FROM gcr.io/kaniko-project/executor:v1.6.0-debug
LABEL maintainer="Codestation <codestation404@gmail.com>"

COPY --from=builder /src/release/drone-kaniko /kaniko/

ENTRYPOINT ["/kaniko/drone-kaniko"]
