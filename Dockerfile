FROM golang:1.14-alpine as builder

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

FROM gcr.io/kaniko-project/executor:debug-23e3fe748663c191bb4e00aa3ac72cd4d22e608c

COPY --from=builder /src/release/drone-kaniko /drone-kaniko

ENTRYPOINT ["/drone-kaniko"]
