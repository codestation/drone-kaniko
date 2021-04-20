FROM golang:1.16-alpine as builder

ARG DRONE_COMMIT_SHA=dev
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
ENV CGO_ENABLED 0
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/
RUN go build -o release/drone-kaniko \
   -ldflags "-w -s \
   -X main.version=${DRONE_COMMIT_SHA:0:8}" \
   -tags netgo \
   ./cmd/drone-kaniko

# use the debug image since it comes with /kaniko/warmer
FROM gcr.io/kaniko-project/executor:v1.5.2-debug

COPY --from=builder /src/release/drone-kaniko /kaniko/

ENTRYPOINT ["/kaniko/drone-kaniko"]
