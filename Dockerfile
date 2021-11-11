FROM golang:1.17-alpine as builder

ARG CI_COMMIT_TAG
ARG CI_COMMIT_BRANCH
ARG CI_COMMIT_SHA
ARG CI_PIPELINE_CREATED_AT
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/
RUN CGO_ENABLED=0 go build -o release/drone-kaniko \
    -ldflags "-w -s \
   -X main.Version=${CI_COMMIT_TAG:-$CI_COMMIT_BRANCH} \
   -X main.Commit=${CI_COMMIT_SHA:0:8} \
   -X main.BuildTime=${CI_PIPELINE_CREATED_AT}" \
    -tags netgo \
    ./cmd/drone-kaniko

# use the debug image since it comes with /kaniko/warmer
FROM gcr.io/kaniko-project/executor:v1.7.0-debug
LABEL maintainer="Codestation <codestation404@gmail.com>"

COPY --from=builder /src/release/drone-kaniko /kaniko/

ENTRYPOINT ["/kaniko/drone-kaniko"]
