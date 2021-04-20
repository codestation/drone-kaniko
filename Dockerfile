FROM golang:1.16-alpine as builder

ARG CI_COMMIT_TAG
ARG CI_COMMIT_BRANCH
ARG CI_COMMIT_SHA
ARG CI_PIPELINE_CREATED_AT
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
ENV CGO_ENABLED 0
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/
RUN go build -o release/drone-kaniko \
    -ldflags "-w -s \
    -X main.version=${CI_COMMIT_TAG:-$CI_COMMIT_BRANCH} \
    -X main.commit=${CI_COMMIT_SHA:0:8} \
    -X main.buildTime=${CI_PIPELINE_CREATED_AT}" \
    -tags netgo \
    ./cmd/drone-kaniko

# use the debug image since it comes with /kaniko/warmer
FROM gcr.io/kaniko-project/executor:v1.5.2-debug

COPY --from=builder /src/release/drone-kaniko /kaniko/

ENTRYPOINT ["/kaniko/drone-kaniko"]
